// Typed API client. Handles CSRF token fetch and injection automatically.

let csrfToken: string | null = null

async function fetchCSRF(): Promise<string> {
  const res = await fetch('/api/csrf-token', { credentials: 'include' })
  if (!res.ok) throw new Error('unauthenticated')
  const data = await res.json()
  return data.csrf_token as string
}

async function getCSRF(): Promise<string> {
  if (!csrfToken) {
    csrfToken = await fetchCSRF()
  }
  return csrfToken
}

export function clearCSRF() {
  csrfToken = null
}

async function request<T>(
  method: string,
  path: string,
  body?: unknown,
  withCSRF = false,
): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }

  if (withCSRF) {
    headers['X-CSRF-Token'] = await getCSRF()
  }

  const res = await fetch(path, {
    method,
    headers,
    credentials: 'include',
    body: body !== undefined ? JSON.stringify(body) : undefined,
  })

  if (res.status === 401) {
    csrfToken = null
    throw new Error('unauthenticated')
  }

  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || `HTTP ${res.status}`)
  }

  if (res.status === 204 || res.headers.get('Content-Length') === '0') {
    return undefined as T
  }

  return res.json() as Promise<T>
}

export const api = {
  login: (email: string, password: string) =>
    request<void>('POST', '/api/login', { email, password }),

  logout: () => request<void>('POST', '/api/logout', undefined, true),

  checkSession: () => fetchCSRF(),
}
