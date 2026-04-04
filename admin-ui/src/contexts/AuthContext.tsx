import { createContext, useContext, useEffect, useState } from 'react'
import type { ReactNode } from 'react'
import { checkAuth, logout as doLogout } from '../lib/auth'
import type { AuthState } from '../lib/auth'

interface AuthContextValue {
  authState: AuthState
  login: () => void
  logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [authState, setAuthState] = useState<AuthState>('loading')

  useEffect(() => {
    checkAuth().then((ok) =>
      setAuthState(ok ? 'authenticated' : 'unauthenticated'),
    )
  }, [])

  function login() {
    setAuthState('authenticated')
  }

  async function logout() {
    await doLogout()
    setAuthState('unauthenticated')
  }

  return (
    <AuthContext.Provider value={{ authState, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

// useAuth is a hook, not a component — co-locating with AuthProvider is intentional
// eslint-disable-next-line react-refresh/only-export-components
export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
