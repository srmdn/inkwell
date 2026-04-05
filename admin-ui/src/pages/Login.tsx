import { useState, useEffect } from 'react'
import type { FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../lib/api'
import type { DemoInfo } from '../lib/api'
import { useAuth } from '../contexts/AuthContext'

export function Login() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [demoInfo, setDemoInfo] = useState<DemoInfo | null>(null)
  const navigate = useNavigate()
  const { login } = useAuth()

  useEffect(() => {
    api.getDemo().then(setDemoInfo).catch(() => {})
  }, [])

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await api.login(email, password)
      login()
      navigate('/admin/posts', { replace: true })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed')
    } finally {
      setLoading(false)
    }
  }

  function fillDemo() {
    if (demoInfo?.email && demoInfo?.password) {
      setEmail(demoInfo.email)
      setPassword(demoInfo.password)
    }
  }

  return (
    <div className="login-container">
      <div className="login-card">
        <h1 className="login-title">FolioCMS</h1>
        <p className="login-subtitle">Admin Dashboard</p>

        {demoInfo?.demo && (
          <div className="demo-banner">
            <p className="demo-banner-label">Demo instance — explore freely, reset anytime.</p>
            <div className="demo-banner-creds">
              <span><strong>Email:</strong> {demoInfo.email}</span>
              <span><strong>Password:</strong> {demoInfo.password}</span>
            </div>
            <button type="button" className="demo-banner-fill" onClick={fillDemo}>
              Fill credentials
            </button>
          </div>
        )}

        <form onSubmit={handleSubmit} className="login-form">
          <div className="field">
            <label htmlFor="email">Email</label>
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              autoFocus
              autoComplete="email"
              disabled={loading}
            />
          </div>

          <div className="field">
            <label htmlFor="password">Password</label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              autoComplete="current-password"
              disabled={loading}
            />
          </div>

          {error && <p className="error-msg">{error}</p>}

          <button type="submit" className="btn-primary" disabled={loading}>
            {loading ? 'Signing in...' : 'Sign in'}
          </button>
        </form>
      </div>
    </div>
  )
}
