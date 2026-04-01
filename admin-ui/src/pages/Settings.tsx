import { useState, useEffect, useRef } from 'react'
import { api } from '../lib/api'
import type { RebuildStatus } from '../lib/api'

function formatTime(ts: string | undefined): string {
  if (!ts || ts.startsWith('0001')) return ''
  return new Date(ts).toLocaleString()
}

export function Settings() {
  const [status, setStatus] = useState<RebuildStatus | null>(null)
  const [triggering, setTriggering] = useState(false)
  const [loadError, setLoadError] = useState('')
  const [triggerError, setTriggerError] = useState('')
  const pollRef = useRef<ReturnType<typeof setInterval> | null>(null)

  useEffect(() => {
    loadStatus()
    return () => stopPolling()
  }, [])

  function stopPolling() {
    if (pollRef.current !== null) {
      clearInterval(pollRef.current)
      pollRef.current = null
    }
  }

  function startPolling() {
    stopPolling()
    pollRef.current = setInterval(async () => {
      try {
        const s = await api.getRebuildStatus()
        setStatus(s)
        if (s.status !== 'running') stopPolling()
      } catch {
        // silently keep polling; network blip should not stop the poller
      }
    }, 2000)
  }

  async function loadStatus() {
    try {
      const s = await api.getRebuildStatus()
      setStatus(s)
      if (s.status === 'running') startPolling()
    } catch (err) {
      setLoadError(err instanceof Error ? err.message : 'Failed to load rebuild status')
    }
  }

  async function handleRebuild() {
    setTriggering(true)
    setTriggerError('')
    try {
      await api.triggerRebuild()
      setStatus((prev) =>
        prev ? { ...prev, status: 'running', output: '', error: '' } : null,
      )
      startPolling()
    } catch (err) {
      setTriggerError(err instanceof Error ? err.message : 'Failed to trigger rebuild')
    } finally {
      setTriggering(false)
    }
  }

  const isRunning = status?.status === 'running' || triggering
  const started = formatTime(status?.started_at)
  const finished = formatTime(status?.finished_at)

  return (
    <>
      <div className="page-header">
        <h1>Settings</h1>
      </div>

      <div className="page-content">
        {loadError && <p className="state-error">{loadError}</p>}

        <div className="settings-card">
          <div className="settings-card-header">
            <div>
              <div className="settings-card-title">Rebuild Site</div>
              <div className="settings-card-desc">
                Triggers a full theme build and restarts the frontend service.
              </div>
            </div>
            <button
              className="btn-sm btn-sm-primary"
              onClick={handleRebuild}
              disabled={isRunning}
            >
              {isRunning ? 'Building...' : 'Rebuild'}
            </button>
          </div>

          {triggerError && <p className="state-error" style={{ margin: '1rem 0 0' }}>{triggerError}</p>}

          {status && status.status !== 'idle' && (
            <div className="rebuild-status">
              <div className="rebuild-status-row">
                <span className="rebuild-label">Status</span>
                <span className={`rebuild-badge rebuild-badge-${status.status}`}>
                  {status.status === 'running' && (
                    <span className="rebuild-dot-spinner" />
                  )}
                  {status.status}
                </span>
              </div>

              {started && (
                <div className="rebuild-status-row">
                  <span className="rebuild-label">Started</span>
                  <span className="rebuild-value">{started}</span>
                </div>
              )}

              {finished && (
                <div className="rebuild-status-row">
                  <span className="rebuild-label">Finished</span>
                  <span className="rebuild-value">{finished}</span>
                </div>
              )}

              {status.error && (
                <div className="rebuild-error">{status.error}</div>
              )}

              {status.output && (
                <pre className="rebuild-output">{status.output}</pre>
              )}
            </div>
          )}
        </div>
      </div>
    </>
  )
}
