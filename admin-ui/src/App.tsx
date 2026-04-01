import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from './contexts/AuthContext'
import { ProtectedRoute } from './components/ProtectedRoute'
import { Layout } from './components/Layout'
import { Login } from './pages/Login'
import { Posts } from './pages/Posts'
import { PostEditor } from './pages/PostEditor'
import { Subscribers } from './pages/Subscribers'
import { Settings } from './pages/Settings'
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
                <Layout>
                  <Posts />
                </Layout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/admin/posts/new"
            element={
              <ProtectedRoute>
                <Layout>
                  <PostEditor />
                </Layout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/admin/posts/:slug/edit"
            element={
              <ProtectedRoute>
                <Layout>
                  <PostEditor />
                </Layout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/admin/subscribers"
            element={
              <ProtectedRoute>
                <Layout>
                  <Subscribers />
                </Layout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/admin/settings"
            element={
              <ProtectedRoute>
                <Layout>
                  <Settings />
                </Layout>
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
