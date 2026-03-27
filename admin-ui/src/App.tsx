import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from './contexts/AuthContext'
import { ProtectedRoute } from './components/ProtectedRoute'
import { Login } from './pages/Login'
import { Posts } from './pages/Posts'
import './index.css'

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/admin/login" element={<Login />} />
          <Route
            path="/admin/posts"
            element={
              <ProtectedRoute>
                <Posts />
              </ProtectedRoute>
            }
          />
          <Route path="/admin" element={<Navigate to="/admin/posts" replace />} />
          <Route path="/admin/*" element={<Navigate to="/admin/posts" replace />} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  )
}
