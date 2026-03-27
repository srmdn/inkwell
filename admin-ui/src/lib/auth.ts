// Auth state helpers. Session is stored server-side as an HttpOnly cookie.
// The SPA detects auth state by calling GET /api/csrf-token (401 = unauthenticated).

import { api, clearCSRF } from './api'

export type AuthState = 'loading' | 'authenticated' | 'unauthenticated'

export async function checkAuth(): Promise<boolean> {
  try {
    await api.checkSession()
    return true
  } catch {
    return false
  }
}

export async function logout(): Promise<void> {
  try {
    await api.logout()
  } finally {
    clearCSRF()
  }
}
