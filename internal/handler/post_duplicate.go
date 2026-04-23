package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/srmdn/foliocms/internal/db"
	"github.com/srmdn/foliocms/internal/model"
	"github.com/srmdn/foliocms/internal/storage"
)

// DuplicatePost copies an existing post as a new draft.
// The new post gets title "Copy of <original>", an auto-generated slug, and draft=true.
func (h *Handler) DuplicatePost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var src model.Post
	err := scanPost(h.db.QueryRow(
		`SELECT id, slug, title, description, tags, draft, publish_date, created_at, updated_at
		 FROM posts WHERE slug = ?`, slug,
	), &src)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, http.StatusNotFound, "post not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not fetch post")
		return
	}

	store := storage.New(h.cfg.ContentDir)
	pf, err := store.Read(slug)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not read post content")
		return
	}

	newSlug := uniquePostSlug(store, h.db, "copy-of-"+slug)
	newTitle := "Copy of " + src.Title

	if err := store.Copy(slug, newSlug); err != nil {
		writeError(w, http.StatusInternalServerError, "could not copy post files")
		return
	}

	today := storage.FormatPublishDate(time.Now())
	if err := store.Write(newSlug, &storage.PostFile{
		Frontmatter: storage.Frontmatter{
			Title:       newTitle,
			Description: pf.Frontmatter.Description,
			PublishDate: today,
			Draft:       true,
			Tags:        pf.Frontmatter.Tags,
			HeroImage:   pf.Frontmatter.HeroImage,
		},
		Body: pf.Body,
	}); err != nil {
		store.Delete(newSlug)
		writeError(w, http.StatusInternalServerError, "could not write post file")
		return
	}

	_, err = h.db.Exec(
		`INSERT INTO posts (slug, title, description, tags, draft, publish_date)
		 VALUES (?, ?, ?, ?, 1, ?)`,
		newSlug, newTitle, src.Description,
		strings.Join(pf.Frontmatter.Tags, ","), today,
	)
	if err != nil {
		store.Delete(newSlug)
		writeError(w, http.StatusInternalServerError, "could not save post")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"slug": newSlug})
}

// uniquePostSlug returns base if it is not taken, otherwise base-2, base-3, ...
func uniquePostSlug(store *storage.Storage, database *db.DB, base string) string {
	candidate := base
	for i := 2; ; i++ {
		var count int
		database.QueryRow(`SELECT COUNT(*) FROM posts WHERE slug = ?`, candidate).Scan(&count)
		if count == 0 && !store.Exists(candidate) {
			return candidate
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
}
