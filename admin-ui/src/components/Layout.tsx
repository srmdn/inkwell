import { NavLink, useNavigate, useLocation } from 'react-router-dom'
import { useState, useEffect, useCallback } from 'react'
import type { ReactNode } from 'react'
import { useAuth } from '../contexts/AuthContext'

const SIDEBAR_KEY = 'folio-sidebar-collapsed'
const MOBILE_BREAKPOINT = 1024

function isMobile() {
  return window.innerWidth < MOBILE_BREAKPOINT
}

export function Layout({ children }: { children: ReactNode }) {
  const { logout } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()

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
    if (mobile) setCollapsed(true)
  }, [location.pathname, mobile])

  const toggle = useCallback(() => {
    setCollapsed(prev => {
      const next = !prev
      if (!mobile) localStorage.setItem(SIDEBAR_KEY, String(next))
      return next
    })
  }, [mobile])

  async function handleLogout() {
    await logout()
    navigate('/admin/login', { replace: true })
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
          <span
            className="sidebar-nav-link disabled"
            title={collapsed && !mobile ? 'Settings' : undefined}
          >
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" aria-hidden="true">
              <circle cx="12" cy="12" r="3" />
              <path d="M19.07 4.93a10 10 0 0 1 0 14.14M4.93 4.93a10 10 0 0 0 0 14.14" />
            </svg>
            {sidebarOpen && <span>Settings</span>}
          </span>
        </nav>

        <div className="sidebar-footer">
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
