import { useState, useEffect, useRef, useCallback } from 'react'
import { api } from '../lib/api'
import type { RebuildStatus, SiteSettings } from '../lib/api'

function formatTime(ts: string | undefined): string {
  if (!ts || ts.startsWith('0001')) return ''
  return new Date(ts).toLocaleString()
}

const defaultSettings: SiteSettings = {
  site_name: '',
  site_description: '',
  social_github: '',
  social_twitter: '',
  social_linkedin: '',
}

export function Settings() {
  // Site settings state
  const [settings, setSettings] = useState<SiteSettings>(defaultSettings)
  const [settingsLoading, setSettingsLoading] = useState(true)
  const [settingsLoadError, setSettingsLoadError] = useState('')
  const [saving, setSaving] = useState(false)
  const [saveError, setSaveError] = useState('')
  const [saveSuccess, setSaveSuccess] = useState(false)

  // Rebuild state
  const [status, setStatus] = useState<RebuildStatus | null>(null)
  const [triggering, setTriggering] = useState(false)
  const [loadError, setLoadError] = useState('')
  const [triggerError, setTriggerError] = useState('')
  const pollRef = useRef<ReturnType<typeof setInterval> | null>(null)

  useEffect(() => {
    loadSettings()
    loadStatus()
    return () => stopPolling()
  }, [])

  // Cmd/Ctrl+S to save settings
  const handleSave = useCallback(async () => {
    setSaving(true)
    setSaveError('')
    setSaveSuccess(false)
    try {
      await api.updateSettings(settings)
      setSaveSuccess(true)
      setTimeout(() => setSaveSuccess(false), 3000)
    } catch (err) {
      setSaveError(err instanceof Error ? err.message : 'Failed to save settings')
    } finally {
      setSaving(false)
    }
  }, [settings])

  useEffect(() => {
    function onKeyDown(e: KeyboardEvent) {
      if ((e.metaKey || e.ctrlKey) && e.key === 's') {
        e.preventDefault()
        handleSave()
      }
    }
    window.addEventListener('keydown', onKeyDown)
    return () => window.removeEventListener('keydown', onKeyDown)
  }, [handleSave])

  async function loadSettings() {
    try {
      const s = await api.getSettings()
      setSettings({
        site_name: s.site_name ?? '',
        site_description: s.site_description ?? '',
        social_github: s.social_github ?? '',
        social_twitter: s.social_twitter ?? '',
        social_linkedin: s.social_linkedin ?? '',
      })
    } catch (err) {
      setSettingsLoadError(err instanceof Error ? err.message : 'Failed to load settings')
    } finally {
      setSettingsLoading(false)
    }
  }

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
        {/* Site Settings */}
        <div className="settings-card" style={{ marginBottom: '1.5rem' }}>
          <div className="settings-card-header">
            <div>
              <div className="settings-card-title">Site Settings</div>
              <div className="settings-card-desc">
                Basic information about your site. Themes read these via{' '}
                <code>/api/settings</code>.
              </div>
            </div>
            <button
              className="btn-sm btn-sm-primary"
              onClick={handleSave}
              disabled={saving || settingsLoading}
            >
              {saving ? 'Saving...' : 'Save'}
            </button>
          </div>

          {settingsLoadError && (
            <p className="state-error" style={{ margin: '1rem 0 0' }}>{settingsLoadError}</p>
          )}
          {saveError && (
            <p className="state-error" style={{ margin: '1rem 0 0' }}>{saveError}</p>
          )}
          {saveSuccess && (
            <p className="state-success" style={{ margin: '1rem 0 0' }}>Settings saved.</p>
          )}

          {!settingsLoading && !settingsLoadError && (
            <div className="settings-form">
              <div className="field">
                <label htmlFor="site_name">Site Name</label>
                <input
                  id="site_name"
                  type="text"
                  value={settings.site_name}
                  onChange={(e) => setSettings((s) => ({ ...s, site_name: e.target.value }))}
                  placeholder="My Site"
                />
              </div>

              <div className="field">
                <label htmlFor="site_description">Description</label>
                <textarea
                  id="site_description"
                  rows={2}
                  value={settings.site_description}
                  onChange={(e) => setSettings((s) => ({ ...s, site_description: e.target.value }))}
                  placeholder="A short description of your site"
                />
              </div>

              <div className="field-group-label">Social Links</div>

              <div className="field">
                <label htmlFor="social_github">GitHub</label>
                <input
                  id="social_github"
                  type="url"
                  value={settings.social_github}
                  onChange={(e) => setSettings((s) => ({ ...s, social_github: e.target.value }))}
                  placeholder="https://github.com/yourusername"
                />
              </div>

              <div className="field">
                <label htmlFor="social_twitter">X / Twitter</label>
                <input
                  id="social_twitter"
                  type="url"
                  value={settings.social_twitter}
                  onChange={(e) => setSettings((s) => ({ ...s, social_twitter: e.target.value }))}
                  placeholder="https://x.com/yourusername"
                />
              </div>

              <div className="field">
                <label htmlFor="social_linkedin">LinkedIn</label>
                <input
                  id="social_linkedin"
                  type="url"
                  value={settings.social_linkedin}
                  onChange={(e) => setSettings((s) => ({ ...s, social_linkedin: e.target.value }))}
                  placeholder="https://linkedin.com/in/yourusername"
                />
              </div>
            </div>
          )}
        </div>

        {/* Rebuild Site */}
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
