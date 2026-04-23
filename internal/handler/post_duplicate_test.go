package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/srmdn/foliocms/internal/config"
	"github.com/srmdn/foliocms/internal/handler"
)

func routerWithDuplicate(h *handler.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/api/admin/posts/{slug}", h.CreatePost)
	r.Post("/api/admin/posts/{slug}/duplicate", h.DuplicatePost)
	r.Get("/api/admin/posts/{slug}", h.GetAdminPost)
	return r
}

func newDuplicateHandler(t *testing.T, contentDir string) *handler.Handler {
	t.Helper()
	database := newTestDB(t)
	cfg := &config.Config{
		ContentDir: contentDir,
		SiteURL:    "https://example.com",
	}
	return handler.New(database, cfg)
}

func duplicatePost(t *testing.T, r *chi.Mux, slug string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/admin/posts/"+slug+"/duplicate", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestDuplicatePostSuccess(t *testing.T) {
	dir := t.TempDir()
	h := newDuplicateHandler(t, dir)
	r := routerWithDuplicate(h)

	createPost(t, r, "my-post", map[string]any{
		"title":       "My Post",
		"description": "A description",
		"draft":       false,
		"tags":        []string{"go", "linux"},
		"body":        "# Hello\n\nSome content.",
	})

	w := duplicatePost(t, r, "my-post")
	if w.Code != http.StatusCreated {
		t.Fatalf("status: got %d, want %d\nbody: %s", w.Code, http.StatusCreated, w.Body.String())
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	newSlug := resp["slug"]
	if newSlug == "" {
		t.Fatalf("response missing slug field")
	}
	if newSlug == "my-post" {
		t.Errorf("duplicate slug must differ from original")
	}

	// Fetch the duplicate and verify it is a draft with "Copy of" title.
	req := httptest.NewRequest(http.MethodGet, "/api/admin/posts/"+newSlug, nil)
	wr := httptest.NewRecorder()
	r.ServeHTTP(wr, req)
	if wr.Code != http.StatusOK {
		t.Fatalf("fetch duplicate: status %d, body %s", wr.Code, wr.Body.String())
	}

	var post map[string]any
	json.Unmarshal(wr.Body.Bytes(), &post)
	if !strings.HasPrefix(post["title"].(string), "Copy of") {
		t.Errorf("title: got %q, want prefix 'Copy of'", post["title"])
	}
	if post["draft"] != true {
		t.Errorf("duplicate must be a draft")
	}
}

func TestDuplicatePostNotFound(t *testing.T) {
	dir := t.TempDir()
	h := newDuplicateHandler(t, dir)
	r := routerWithDuplicate(h)

	w := duplicatePost(t, r, "no-such-post")
	if w.Code != http.StatusNotFound {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestDuplicatePostUniqueSlug(t *testing.T) {
	dir := t.TempDir()
	h := newDuplicateHandler(t, dir)
	r := routerWithDuplicate(h)

	createPost(t, r, "article", map[string]any{"title": "Article", "draft": false})

	w1 := duplicatePost(t, r, "article")
	w2 := duplicatePost(t, r, "article")

	if w1.Code != http.StatusCreated || w2.Code != http.StatusCreated {
		t.Fatalf("both duplicates must succeed: %d, %d", w1.Code, w2.Code)
	}

	var r1, r2 map[string]string
	json.Unmarshal(w1.Body.Bytes(), &r1)
	json.Unmarshal(w2.Body.Bytes(), &r2)

	if r1["slug"] == r2["slug"] {
		t.Errorf("two duplicates must get different slugs, both got %q", r1["slug"])
	}
}

func TestDuplicatePostPreservesBody(t *testing.T) {
	dir := t.TempDir()
	h := newDuplicateHandler(t, dir)
	r := routerWithDuplicate(h)

	createPost(t, r, "rich-post", map[string]any{
		"title": "Rich Post",
		"draft": false,
		"body":  "## Section\n\nContent here.",
	})

	w := duplicatePost(t, r, "rich-post")
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/posts/"+resp["slug"], nil)
	wr := httptest.NewRecorder()
	r.ServeHTTP(wr, req)

	var post map[string]any
	json.Unmarshal(wr.Body.Bytes(), &post)
	if !strings.Contains(post["body"].(string), "Content here.") {
		t.Errorf("duplicate must preserve body content")
	}
}
