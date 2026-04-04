import { useEffect, useRef, useState } from 'react'
import { api } from '../lib/api'
import type { MediaFile } from '../lib/api'
import { MediaGrid } from '../components/media/MediaGrid'

const ACCEPTED = ['image/jpeg', 'image/png', 'image/gif', 'image/webp', 'image/svg+xml']

export function Media() {
  const [files, setFiles] = useState<MediaFile[] | null>(null)
  const [error, setError] = useState('')
  const [uploading, setUploading] = useState(false)
  const [dragOver, setDragOver] = useState(false)
  const [copied, setCopied] = useState<string | null>(null)
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    load()
  }, [])

  async function load() {
    try {
      const data = await api.listMedia()
      setFiles(data ?? [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load media')
    }
  }

  async function upload(file: File) {
    if (!ACCEPTED.includes(file.type)) {
      setError('Only image files are allowed (JPEG, PNG, GIF, WebP, SVG).')
      return
    }
    setError('')
    setUploading(true)
    try {
      const mf = await api.uploadMedia(file)
      setFiles((prev) => (prev ? [mf, ...prev] : [mf]))
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Upload failed')
    } finally {
      setUploading(false)
    }
  }

  function handleFileInput(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0]
    if (file) upload(file)
    e.target.value = ''
  }

  function handleDrop(e: React.DragEvent) {
    e.preventDefault()
    setDragOver(false)
    const file = e.dataTransfer.files?.[0]
    if (file) upload(file)
  }

  function handleDragOver(e: React.DragEvent) {
    e.preventDefault()
    setDragOver(true)
  }

  function handleDragLeave() {
    setDragOver(false)
  }

  async function handleDelete(key: string) {
    await api.deleteMedia(key)
    setFiles((prev) => prev?.filter((f) => f.key !== key) ?? [])
  }

  function handleCopy(url: string, key: string) {
    navigator.clipboard.writeText(url).then(() => {
      setCopied(key)
      setTimeout(() => setCopied(null), 1800)
    })
  }

  const count = files?.length ?? 0

  return (
    <>
      <div className="page-header">
        <h1>Media</h1>
        {files !== null && (
          <span className="page-count">
            {count} {count === 1 ? 'file' : 'files'}
          </span>
        )}
      </div>

      <div className="page-content">
        {error && <p className="state-error">{error}</p>}

        <div
          className={`upload-zone${dragOver ? ' upload-zone-active' : ''}`}
          onDrop={handleDrop}
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
        >
          <input
            ref={inputRef}
            type="file"
            accept={ACCEPTED.join(',')}
            style={{ display: 'none' }}
            onChange={handleFileInput}
          />
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" aria-hidden="true">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
            <polyline points="17 8 12 3 7 8" />
            <line x1="12" y1="3" x2="12" y2="15" />
          </svg>
          <span className="upload-zone-text">
            {uploading ? 'Uploading...' : 'Drop image here or'}
          </span>
          {!uploading && (
            <button
              className="btn-sm btn-sm-primary"
              onClick={() => inputRef.current?.click()}
            >
              Choose file
            </button>
          )}
        </div>

        {files === null && !error && (
          <div className="state-center">
            <div className="spinner" />
          </div>
        )}

        {files !== null && (
          <MediaGrid
            files={files}
            onDelete={handleDelete}
            copied={copied}
            onCopy={handleCopy}
          />
        )}
      </div>
    </>
  )
}
