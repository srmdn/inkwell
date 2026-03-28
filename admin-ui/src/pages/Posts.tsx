import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { api } from '../lib/api'
import type { Post } from '../lib/api'

type PostStatus = 'published' | 'draft' | 'scheduled'

function getStatus(post: Post): PostStatus {
  if (post.draft) return 'draft'
  if (new Date(post.publish_date) > new Date()) return 'scheduled'
  return 'published'
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

export function Posts() {
  const [posts, setPosts] = useState<Post[] | null>(null)
  const [error, setError] = useState('')
  const [confirmSlug, setConfirmSlug] = useState<string | null>(null)
  const [deleting, setDeleting] = useState(false)

  useEffect(() => {
    loadPosts()
  }, [])

  async function loadPosts() {
    try {
      const data = await api.getPosts()
      setPosts(data ?? [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load posts')
    }
  }

  async function handleDelete(slug: string) {
    setDeleting(true)
    try {
      await api.deletePost(slug)
      setPosts((prev) => prev?.filter((p) => p.slug !== slug) ?? [])
      setConfirmSlug(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete post')
    } finally {
      setDeleting(false)
    }
  }

  return (
    <>
      <div className="page-header">
        <h1>Posts</h1>
        <Link to="/admin/posts/new" className="btn-sm btn-sm-primary">
          + New post
        </Link>
      </div>

      <div className="page-content">
        {error && <p className="state-error">{error}</p>}

        {posts === null && !error && (
          <div className="state-center">
            <div className="spinner" />
          </div>
        )}

        {posts !== null && posts.length === 0 && (
          <div className="state-center state-muted">
            No posts yet.{' '}
            <Link to="/admin/posts/new">Create your first post.</Link>
          </div>
        )}

        {posts !== null && posts.length > 0 && (
          <div className="post-list">
            {posts.map((post) => {
              const status = getStatus(post)
              return (
                <div key={post.slug} className="post-item">
                  <div className="post-info">
                    <div className="post-title">{post.title}</div>
                    <div className="post-sub">
                      <span className={`badge badge-${status}`}>{status}</span>
                      <span className="post-date">{formatDate(post.publish_date)}</span>
                    </div>
                  </div>

                  <div className="post-actions">
                    {confirmSlug === post.slug ? (
                      <div className="confirm-inline">
                        <span>Delete?</span>
                        <button
                          className="btn-sm btn-sm-danger"
                          onClick={() => handleDelete(post.slug)}
                          disabled={deleting}
                        >
                          {deleting ? '...' : 'Yes'}
                        </button>
                        <button
                          className="btn-sm btn-sm-ghost"
                          onClick={() => setConfirmSlug(null)}
                          disabled={deleting}
                        >
                          Cancel
                        </button>
                      </div>
                    ) : (
                      <>
                        <Link
                          to={`/admin/posts/${post.slug}/edit`}
                          className="btn-sm btn-sm-ghost"
                        >
                          Edit
                        </Link>
                        <button
                          className="btn-sm btn-sm-danger"
                          onClick={() => setConfirmSlug(post.slug)}
                        >
                          Delete
                        </button>
                      </>
                    )}
                  </div>
                </div>
              )
            })}
          </div>
        )}
      </div>
    </>
  )
}
