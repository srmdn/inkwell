import { useEffect, useRef } from 'react'
import { Editor, rootCtx, defaultValueCtx } from '@milkdown/core'
import { commonmark } from '@milkdown/preset-commonmark'
import { gfm } from '@milkdown/preset-gfm'
import { history, undoCommand, redoCommand } from '@milkdown/plugin-history'
import { listener, listenerCtx } from '@milkdown/plugin-listener'
import { callCommand } from '@milkdown/utils'
import {
  toggleEmphasisCommand,
  toggleStrongCommand,
  toggleInlineCodeCommand,
  wrapInBlockquoteCommand,
  wrapInBulletListCommand,
  wrapInOrderedListCommand,
  insertHrCommand,
  turnIntoTextCommand,
} from '@milkdown/preset-commonmark'
import { insertTableCommand } from '@milkdown/preset-gfm'

interface EditorProps {
  value: string
  onChange: (markdown: string) => void
}

export function MarkdownEditor({ value, onChange }: EditorProps) {
  const editorRef = useRef<HTMLDivElement>(null)
  const instanceRef = useRef<Editor | null>(null)
  const onChangeRef = useRef(onChange)
  onChangeRef.current = onChange

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
  }, []) // mount once — value is set via defaultValueCtx

  function action(cmd: Parameters<typeof callCommand>[0]) {
    instanceRef.current?.action(callCommand(cmd))
  }

  return (
    <div className="editor-wrap">
      <div className="editor-toolbar">
        <button type="button" className="tb-btn" title="Bold" onClick={() => action(toggleStrongCommand.key)}>
          <b>B</b>
        </button>
        <button type="button" className="tb-btn" title="Italic" onClick={() => action(toggleEmphasisCommand.key)}>
          <i>I</i>
        </button>
        <button type="button" className="tb-btn" title="Inline code" onClick={() => action(toggleInlineCodeCommand.key)}>
          {'</>'}
        </button>
        <span className="tb-sep" />
        <button type="button" className="tb-btn" title="Blockquote" onClick={() => action(wrapInBlockquoteCommand.key)}>
          &ldquo;
        </button>
        <button type="button" className="tb-btn" title="Bullet list" onClick={() => action(wrapInBulletListCommand.key)}>
          &#8226;&#8212;
        </button>
        <button type="button" className="tb-btn" title="Ordered list" onClick={() => action(wrapInOrderedListCommand.key)}>
          1&#8212;
        </button>
        <span className="tb-sep" />
        <button type="button" className="tb-btn" title="Table" onClick={() => action(insertTableCommand.key)}>
          &#9956;
        </button>
        <button type="button" className="tb-btn" title="Horizontal rule" onClick={() => action(insertHrCommand.key)}>
          &#8212;
        </button>
        <button type="button" className="tb-btn" title="Plain text" onClick={() => action(turnIntoTextCommand.key)}>
          T
        </button>
        <span className="tb-sep" />
        <button type="button" className="tb-btn" title="Undo" onClick={() => action(undoCommand.key)}>
          &#8617;
        </button>
        <button type="button" className="tb-btn" title="Redo" onClick={() => action(redoCommand.key)}>
          &#8618;
        </button>
      </div>
      <div ref={editorRef} className="editor-body" />
    </div>
  )
}
