import { useEffect, useState } from 'react'
import { api } from '../lib/api'
import type { Subscriber } from '../lib/api'

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

export function Subscribers() {
  const [subscribers, setSubscribers] = useState<Subscriber[] | null>(null)
  const [error, setError] = useState('')
  const [confirmId, setConfirmId] = useState<number | null>(null)
  const [deleting, setDeleting] = useState(false)

  const [subject, setSubject] = useState('')
  const [body, setBody] = useState('')
  const [sending, setSending] = useState(false)
  const [sendError, setSendError] = useState('')
  const [sendSuccess, setSendSuccess] = useState(false)

  useEffect(() => {
    loadSubscribers()
  }, [])

  async function loadSubscribers() {
    try {
      const data = await api.getSubscribers()
      setSubscribers(data ?? [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load subscribers')
    }
  }

  async function handleDelete(id: number) {
    setDeleting(true)
    try {
      await api.deleteSubscriber(id)
      setSubscribers((prev) => prev?.filter((s) => s.id !== id) ?? [])
      setConfirmId(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to remove subscriber')
    } finally {
      setDeleting(false)
    }
  }

  async function handleSend() {
    setSendError('')
    setSendSuccess(false)
    setSending(true)
    try {
      await api.sendNewsletter(subject, body)
      setSendSuccess(true)
      setSubject('')
      setBody('')
    } catch (err) {
      setSendError(err instanceof Error ? err.message : 'Failed to send newsletter')
    } finally {
      setSending(false)
    }
  }

  const count = subscribers?.length ?? 0

  return (
    <>
      <div className="page-header">
        <h1>Subscribers</h1>
        {subscribers !== null && (
          <span className="page-count">
            {count} {count === 1 ? 'subscriber' : 'subscribers'}
          </span>
        )}
      </div>

      <div className="page-content">
        <div className="settings-card">
          <div className="settings-card-header">
            <div>
              <div className="settings-card-title">Send Newsletter</div>
              <div className="settings-card-desc">Send an email to all subscribers.</div>
            </div>
            <button
              className="btn-sm btn-sm-primary"
              onClick={handleSend}
              disabled={sending || !subject.trim() || !body.trim()}
            >
              {sending ? 'Sending...' : 'Send'}
            </button>
          </div>
          {sendError && <p className="state-error">{sendError}</p>}
          {sendSuccess && <p className="state-success">Newsletter sent.</p>}
          <div className="field">
            <label>Subject</label>
            <input
              type="text"
              value={subject}
              onChange={(e) => setSubject(e.target.value)}
              placeholder="Your newsletter subject"
            />
          </div>
          <div className="field">
            <label>Body</label>
            <textarea
              value={body}
              onChange={(e) => setBody(e.target.value)}
              placeholder="Write your newsletter in plain text..."
              rows={8}
            />
          </div>
        </div>

        {error && <p className="state-error">{error}</p>}

        {subscribers === null && !error && (
          <div className="state-center">
            <div className="spinner" />
          </div>
        )}

        {subscribers !== null && subscribers.length === 0 && (
          <div className="state-center state-muted">No subscribers yet.</div>
        )}

        {subscribers !== null && subscribers.length > 0 && (
          <div className="post-list">
            {subscribers.map((sub) => (
              <div key={sub.id} className="post-item">
                <div className="post-info">
                  <div className="post-title">{sub.email}</div>
                  <div className="post-sub">
                    <span className="post-date">Subscribed {formatDate(sub.subscribed_at)}</span>
                  </div>
                </div>

                <div className="post-actions">
                  {confirmId === sub.id ? (
                    <div className="confirm-inline">
                      <span>Remove?</span>
                      <button
                        className="btn-sm btn-sm-danger"
                        onClick={() => handleDelete(sub.id)}
                        disabled={deleting}
                      >
                        {deleting ? '...' : 'Yes'}
                      </button>
                      <button
                        className="btn-sm btn-sm-ghost"
                        onClick={() => setConfirmId(null)}
                        disabled={deleting}
                      >
                        Cancel
                      </button>
                    </div>
                  ) : (
                    <button
                      className="btn-sm btn-sm-danger"
                      onClick={() => setConfirmId(sub.id)}
                    >
                      Remove
                    </button>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </>
  )
}
