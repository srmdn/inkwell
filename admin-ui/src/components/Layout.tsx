import { NavLink, useNavigate, useLocation } from 'react-router-dom'
import { useState, useEffect, useCallback } from 'react'
import type { ReactNode } from 'react'
import { useAuth } from '../contexts/AuthContext'
import { api } from '../lib/api'

const SIDEBAR_KEY = 'folio-sidebar-collapsed'
const MOBILE_BREAKPOINT = 1024

function isMobile() {
  return window.innerWidth < MOBILE_BREAKPOINT
}

export function Layout({ children }: { children: ReactNode }) {
  const { logout } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()
  const [demoMode, setDemoMode] = useState(false)
  const [resetting, setResetting] = useState(false)

  // Desktop: collapsed = icon-only. Mobile: open = drawer visible.
  const [collapsed, setCollapsed] = useState<boolean>(() => {
    if (typeof window === 'undefined') return false
    if (isMobile()) return true // mobile starts closed
    return localStorage.getItem(SIDEBAR_KEY) === 'true'
  })
  const [mobile, setMobile] = useState<boolean>(() => isMobile())

  // Track mobile/desktop breakpoint
  useEffect(() => {
    function onResize() {
      const nowMobile = isMobile()
      setMobile(nowMobile)
      if (!nowMobile) {
        // Restore desktop preference when resizing back
        setCollapsed(localStorage.getItem(SIDEBAR_KEY) === 'true')
      } else {
        setCollapsed(true) // always close on mobile resize
      }
    }
    window.addEventListener('resize', onResize)
    return () => window.removeEventListener('resize', onResize)
  }, [])

  // Close mobile drawer on route change
  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    if (mobile) setCollapsed(true)
  }, [location.pathname, mobile])

  const toggle = useCallback(() => {
    setCollapsed(prev => {
      const next = !prev
      if (!mobile) localStorage.setItem(SIDEBAR_KEY, String(next))
      return next
    })
  }, [mobile])

  useEffect(() => {
    api.getDemo().then((info) => setDemoMode(info.demo)).catch(() => {})
  }, [])

  async function handleLogout() {
    await logout()
    navigate('/admin/login', { replace: true })
  }

  async function handleReset() {
    if (!confirm('Reset the demo? All changes will be lost.')) return
    setResetting(true)
    try {
      await api.resetDemo()
      navigate('/admin/posts', { replace: true })
      window.location.reload()
    } catch {
      alert('Reset failed. Please try again.')
    } finally {
      setResetting(false)
    }
  }

  const sidebarOpen = !collapsed

  return (
    <div className={`dashboard${mobile ? ' dashboard-mobile' : ''}`}>
      {/* Mobile backdrop */}
      {mobile && sidebarOpen && (
        <div className="sidebar-backdrop" onClick={() => setCollapsed(true)} />
      )}

      <aside className={`sidebar${sidebarOpen ? ' sidebar-open' : ' sidebar-closed'}${collapsed && !mobile ? ' sidebar-collapsed' : ''}`}>
        <div className="sidebar-brand">
          {sidebarOpen && (
            <>
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                <rect width="24" height="24" rx="6" fill="#22c55e" />
                <path d="M7 8h10M7 12h7M7 16h5" stroke="white" strokeWidth="1.8" strokeLinecap="round" />
              </svg>
              <span className="sidebar-brand-text">FolioCMS</span>
            </>
          )}
          {!sidebarOpen && !mobile && (
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" aria-hidden="true">
              <rect width="24" height="24" rx="6" fill="#22c55e" />
              <path d="M7 8h10M7 12h7M7 16h5" stroke="white" strokeWidth="1.8" strokeLinecap="round" />
            </svg>
          )}
          <button
            className="sidebar-toggle"
            onClick={toggle}
            aria-label={sidebarOpen ? 'Collapse sidebar' : 'Expand sidebar'}
          >
            {sidebarOpen ? (
              // chevron-left
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
                <path d="M15 18l-6-6 6-6" />
              </svg>
            ) : (
              // chevron-right / hamburger
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
                <path d="M9 18l6-6-6-6" />
              </svg>
            )}
          </button>
        </div>

        <nav className="sidebar-nav">
          <NavLink
            to="/admin/posts"
            className={({ isActive }) => `sidebar-nav-link${isActive ? ' active' : ''}`}
            title={collapsed && !mobile ? 'Posts' : undefined}
          >
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" aria-hidden="true">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
              <path d="M14 2v6h6M16 13H8M16 17H8M10 9H8" />
            </svg>
            {sidebarOpen && <span>Posts</span>}
          </NavLink>
          <NavLink
            to="/admin/media"
            className={({ isActive }) => `sidebar-nav-link${isActive ? ' active' : ''}`}
            title={collapsed && !mobile ? 'Media' : undefined}
          >
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" aria-hidden="true">
              <rect x="3" y="3" width="18" height="18" rx="2" />
              <circle cx="8.5" cy="8.5" r="1.5" />
              <path d="M21 15l-5-5L5 21" />
            </svg>
            {sidebarOpen && <span>Media</span>}
          </NavLink>
          <NavLink
            to="/admin/subscribers"
            className={({ isActive }) => `sidebar-nav-link${isActive ? ' active' : ''}`}
            title={collapsed && !mobile ? 'Subscribers' : undefined}
          >
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" aria-hidden="true">
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
              <circle cx="9" cy="7" r="4" />
              <path d="M23 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75" />
            </svg>
            {sidebarOpen && <span>Subscribers</span>}
          </NavLink>
          <NavLink
            to="/admin/settings"
            className={({ isActive }) => `sidebar-nav-link${isActive ? ' active' : ''}`}
            title={collapsed && !mobile ? 'Settings' : undefined}
          >
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" aria-hidden="true">
              <circle cx="12" cy="12" r="3" />
              <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z" />
            </svg>
            {sidebarOpen && <span>Settings</span>}
          </NavLink>
        </nav>

        <div className="sidebar-footer">
          {demoMode && (
            <button
              className="sidebar-nav-link sidebar-reset-demo"
              onClick={handleReset}
              disabled={resetting}
              title={collapsed && !mobile ? 'Reset Demo' : undefined}
            >
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" aria-hidden="true">
                <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8" />
                <path d="M3 3v5h5" />
              </svg>
              {sidebarOpen && <span>{resetting ? 'Resetting...' : 'Reset Demo'}</span>}
            </button>
          )}
          <button
            className="sidebar-nav-link"
            onClick={handleLogout}
            title={collapsed && !mobile ? 'Sign out' : undefined}
          >
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" aria-hidden="true">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4M16 17l5-5-5-5M21 12H9" />
            </svg>
            {sidebarOpen && <span>Sign out</span>}
          </button>
        </div>
      </aside>

      <div className={`main-wrapper${collapsed && !mobile ? ' main-wrapper-collapsed' : ''}`}>
        {/* Mobile top bar */}
        {mobile && (
          <div className="mobile-topbar">
            <button className="mobile-menu-btn" onClick={toggle} aria-label="Open menu">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
                <line x1="3" y1="6" x2="21" y2="6" />
                <line x1="3" y1="12" x2="21" y2="12" />
                <line x1="3" y1="18" x2="21" y2="18" />
              </svg>
            </button>
            <div className="mobile-topbar-brand">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                <rect width="24" height="24" rx="6" fill="#22c55e" />
                <path d="M7 8h10M7 12h7M7 16h5" stroke="white" strokeWidth="1.8" strokeLinecap="round" />
              </svg>
              FolioCMS
            </div>
          </div>
        )}
        <main className="main">
          {children}
        </main>
      </div>
    </div>
  )
}
