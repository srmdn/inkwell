import { useEffect, useRef, useState } from 'react'
import { Editor, rootCtx, defaultValueCtx, commandsCtx, editorViewCtx } from '@milkdown/core'
import { commonmark } from '@milkdown/preset-commonmark'
import { gfm } from '@milkdown/preset-gfm'
import { history } from '@milkdown/plugin-history'
import { listener, listenerCtx } from '@milkdown/plugin-listener'
import {
  Bold, Italic, Strikethrough, Code, Code2,
  Heading1, Heading2, Heading3,
  Quote, List, ListOrdered, Minus, Pilcrow, Link, Table,
} from 'lucide-react'

interface EditorProps {
  value: string
  onChange: (markdown: string) => void
}

const CMD = {
  BOLD:         'ToggleStrong',
  ITALIC:       'ToggleEmphasis',
  STRIKE:       'ToggleStrikeThrough',
  CODE:         'ToggleInlineCode',
  HEADING:      'WrapInHeading',
  PARAGRAPH:    'TurnIntoText',
  BLOCKQUOTE:   'WrapInBlockquote',
  BULLET_LIST:  'WrapInBulletList',
  ORDERED_LIST: 'WrapInOrderedList',
  CODE_BLOCK:   'CreateCodeBlock',
  HR:           'InsertHr',
} as const

function TbBtn({
  onClick, title, active, children,
}: {
  onClick: () => void
  title: string
  active?: boolean
  children: React.ReactNode
}) {
  return (
    <button
      type="button"
      title={title}
      onMouseDown={(e) => {
        e.preventDefault() // Prevent editor losing focus
        onClick()
      }}
      className={`tb-btn${active ? ' tb-btn-active' : ''}`}
    >
      {children}
    </button>
  )
}

function Sep() {
  return <span className="tb-sep" />
}

export function MarkdownEditor({ value, onChange }: EditorProps) {
  const editorRef = useRef<HTMLDivElement>(null)
  const toolbarRef = useRef<HTMLDivElement>(null)
  const instanceRef = useRef<Editor | null>(null)
  const onChangeRef = useRef(onChange)
  onChangeRef.current = onChange

  const [activeMarks, setActiveMarks] = useState<Set<string>>(new Set())
  const [activeNode, setActiveNode] = useState<{ type: string; level?: number } | null>(null)

  useEffect(() => {
    if (!editorRef.current) return

    let destroyed = false

    Editor.make()
      .config((ctx) => {
        ctx.set(rootCtx, editorRef.current!)
        ctx.set(defaultValueCtx, value)
        ctx.get(listenerCtx).markdownUpdated((_ctx, markdown) => {
          if (!destroyed) onChangeRef.current(markdown)
        })
      })
      .use(commonmark)
      .use(gfm)
      .use(history)
      .use(listener)
      .create()
      .then((editor) => {
        if (!destroyed) instanceRef.current = editor
        else editor.destroy()
      })

    return () => {
      destroyed = true
      instanceRef.current?.destroy()
      instanceRef.current = null
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  // Sync active state for toolbar highlights
  useEffect(() => {
    let rafId: number
    const loop = () => {
      const editor = instanceRef.current
      if (editor) {
        try {
          editor.action((ctx) => {
            const view = ctx.get(editorViewCtx)
            const { state } = view
            const { schema, selection, storedMarks } = state
            const { $from } = selection

            const marks = storedMarks ?? $from.marks()
            const active = new Set<string>()
            for (const mark of marks) active.add(mark.type.name)
            setActiveMarks(active)

            const node = $from.parent
            if (node.type === schema.nodes.heading) {
              setActiveNode({ type: 'heading', level: node.attrs.level })
            } else if (node.type === schema.nodes.paragraph) {
              setActiveNode({ type: 'paragraph' })
            } else if (
              node.type === schema.nodes.blockquote ||
              $from.node($from.depth - 1)?.type === schema.nodes.blockquote
            ) {
              setActiveNode({ type: 'blockquote' })
            } else {
              setActiveNode({ type: node.type.name })
            }
          })
        } catch {
          // editor not ready yet
        }
      }
      rafId = requestAnimationFrame(loop)
    }
    rafId = requestAnimationFrame(loop)
    return () => cancelAnimationFrame(rafId)
  }, [])


  function callCommand(key: string, payload?: unknown) {
    instanceRef.current?.action((ctx) => {
      ctx.get(commandsCtx).call(key, payload)
    })
  }

  function handleLink() {
    const editor = instanceRef.current
    if (!editor) return
    editor.action((ctx) => {
      const view = ctx.get(editorViewCtx)
      const { state, dispatch } = view
      const { schema, selection } = state
      let { from, to } = selection
      const linkMark = schema.marks.link
      if (!linkMark) return

      if (from === to) {
        const $pos = state.doc.resolve(from)
        if ($pos.marks().some((m: any) => m.type === linkMark)) {
          while (from > 0 && state.doc.resolve(from - 1).marks().some((m: any) => m.type === linkMark)) from--
          while (to < state.doc.content.size && state.doc.resolve(to).marks().some((m: any) => m.type === linkMark)) to++
        }
      }

      if (from === to) {
        window.alert('Select some text first to add a link.')
        view.focus()
        return
      }

      const existing = state.doc.resolve(from).marks().find((m: any) => m.type === linkMark)
      const existingHref = existing?.attrs.href ?? ''
      const url = window.prompt('Link URL (leave empty to remove):', existingHref)
      if (url === null) { view.focus(); return }

      const tr = state.tr
      if (url === '') {
        tr.removeMark(from, to, linkMark)
      } else {
        tr.addMark(from, to, linkMark.create({ href: url }))
      }
      dispatch(tr)
      view.focus()
    })
  }

  function handleTable() {
    const editor = instanceRef.current
    if (!editor) return
    editor.action((ctx) => {
      const view = ctx.get(editorViewCtx)
      const { state, dispatch } = view
      const { schema, selection } = state
      const tableNode = schema.nodes.table
      const rowNode = schema.nodes.table_row
      const cellNode = schema.nodes.table_cell
      const headerNode = schema.nodes.table_header
      if (!tableNode || !rowNode || !cellNode || !headerNode) return

      const cell = (text: string, isHeader = false) => {
        const nodeType = isHeader ? headerNode : cellNode
        return nodeType.create({}, schema.nodes.paragraph.create({}, text ? schema.text(text) : undefined))
      }

      const table = tableNode.create({}, [
        rowNode.create({}, [cell('Header 1', true), cell('Header 2', true), cell('Header 3', true)]),
        rowNode.create({}, [cell(''), cell(''), cell('')]),
        rowNode.create({}, [cell(''), cell(''), cell('')]),
      ])

      const { $from } = selection
      const insertPos = $from.end($from.depth) + 1
      dispatch(state.tr.insert(insertPos, table))
      view.focus()
    })
  }

  return (
    <div className="editor-wrap">
      <div ref={toolbarRef} className="editor-toolbar">
        {/* Block type */}
        <TbBtn onClick={() => callCommand(CMD.PARAGRAPH)} title="Paragraph" active={activeNode?.type === 'paragraph'}>
          <Pilcrow size={14} />
        </TbBtn>
        <TbBtn onClick={() => callCommand(CMD.HEADING, 1)} title="Heading 1" active={activeNode?.type === 'heading' && activeNode.level === 1}>
          <Heading1 size={14} />
        </TbBtn>
        <TbBtn onClick={() => callCommand(CMD.HEADING, 2)} title="Heading 2" active={activeNode?.type === 'heading' && activeNode.level === 2}>
          <Heading2 size={14} />
        </TbBtn>
        <TbBtn onClick={() => callCommand(CMD.HEADING, 3)} title="Heading 3" active={activeNode?.type === 'heading' && activeNode.level === 3}>
          <Heading3 size={14} />
        </TbBtn>

        <Sep />

        {/* Inline marks */}
        <TbBtn onClick={() => callCommand(CMD.BOLD)} title="Bold (Ctrl+B)" active={activeMarks.has('strong')}>
          <Bold size={14} />
        </TbBtn>
        <TbBtn onClick={() => callCommand(CMD.ITALIC)} title="Italic (Ctrl+I)" active={activeMarks.has('em')}>
          <Italic size={14} />
        </TbBtn>
        <TbBtn onClick={() => callCommand(CMD.STRIKE)} title="Strikethrough" active={activeMarks.has('strike_through')}>
          <Strikethrough size={14} />
        </TbBtn>
        <TbBtn onClick={() => callCommand(CMD.CODE)} title="Inline Code" active={activeMarks.has('code_inline')}>
          <Code size={14} />
        </TbBtn>
        <TbBtn onClick={handleLink} title="Link — select text, then click" active={activeMarks.has('link')}>
          <Link size={14} />
        </TbBtn>

        <Sep />

        {/* Block elements */}
        <TbBtn onClick={() => callCommand(CMD.BLOCKQUOTE)} title="Blockquote" active={activeNode?.type === 'blockquote'}>
          <Quote size={14} />
        </TbBtn>
        <TbBtn onClick={() => callCommand(CMD.BULLET_LIST)} title="Bullet List">
          <List size={14} />
        </TbBtn>
        <TbBtn onClick={() => callCommand(CMD.ORDERED_LIST)} title="Numbered List">
          <ListOrdered size={14} />
        </TbBtn>
        <TbBtn onClick={() => callCommand(CMD.CODE_BLOCK)} title="Code Block" active={activeNode?.type === 'code_block'}>
          <Code2 size={14} />
        </TbBtn>
        <TbBtn onClick={handleTable} title="Insert Table">
          <Table size={14} />
        </TbBtn>
        <TbBtn onClick={() => callCommand(CMD.HR)} title="Divider">
          <Minus size={14} />
        </TbBtn>
      </div>

      <div
        ref={editorRef}
        className="editor-body"
        onClick={(e) => {
          if ((e.target as HTMLElement).classList.contains('editor-body')) {
            editorRef.current?.querySelector<HTMLElement>('.ProseMirror')?.focus()
          }
        }}
      />
    </div>
  )
}
