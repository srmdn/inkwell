import { useState } from 'react'
import type { MediaFile } from '../../lib/api'

function formatSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

interface Props {
  files: MediaFile[]
  onDelete: (key: string) => Promise<void>
  copied: string | null
  onCopy: (url: string, key: string) => void
}

export function MediaGrid({ files, onDelete, copied, onCopy }: Props) {
  const [confirmKey, setConfirmKey] = useState<string | null>(null)
  const [deleting, setDeleting] = useState(false)

  async function handleDelete(key: string) {
    setDeleting(true)
    try {
      await onDelete(key)
      setConfirmKey(null)
    } finally {
      setDeleting(false)
    }
  }

  if (files.length === 0) {
    return (
      <div className="state-center state-muted">No media uploaded yet.</div>
    )
  }

  return (
    <div className="media-grid">
      {files.map((f) => (
        <div key={f.key} className="media-card">
          <button
            className="media-thumb-btn"
            onClick={() => onCopy(f.url, f.key)}
            title="Click to copy URL"
          >
            <img src={f.url} alt={f.filename} className="media-thumb" loading="lazy" />
            {copied === f.key && (
              <div className="media-copied-overlay">Copied!</div>
            )}
          </button>
          <div className="media-card-info">
            <div className="media-card-name" title={f.filename}>{f.filename}</div>
            <div className="media-card-meta">
              <span className="media-card-size">{formatSize(f.size)}</span>
              {confirmKey === f.key ? (
                <div className="confirm-inline">
                  <button
                    className="btn-sm btn-sm-danger"
                    onClick={() => handleDelete(f.key)}
                    disabled={deleting}
                  >
                    {deleting ? '...' : 'Delete'}
                  </button>
                  <button
                    className="btn-sm btn-sm-ghost"
                    onClick={() => setConfirmKey(null)}
                    disabled={deleting}
                  >
                    Cancel
                  </button>
                </div>
              ) : (
                <button
                  className="btn-sm btn-sm-danger"
                  onClick={() => setConfirmKey(f.key)}
                >
                  Delete
                </button>
              )}
            </div>
          </div>
        </div>
      ))}
    </div>
  )
}
