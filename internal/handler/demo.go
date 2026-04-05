package handler

import (
	"log"
	"net/http"

	"github.com/srmdn/foliocms/internal/demo"
)

// DemoInfo returns demo credentials (public endpoint, only meaningful when
// DEMO_MODE=true). Returns {demo: false} when demo mode is off so the
// frontend can branch without a 404.
func (h *Handler) DemoInfo(w http.ResponseWriter, r *http.Request) {
	if !h.cfg.DemoMode {
		writeJSON(w, http.StatusOK, map[string]any{"demo": false})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"demo":     true,
		"email":    h.cfg.DemoEmail,
		"password": h.cfg.DemoPasswd,
	})
}

// DemoReset wipes all posts and media and re-applies the seed.
// Only available when DEMO_MODE=true; protected by JWT + CSRF.
func (h *Handler) DemoReset(w http.ResponseWriter, r *http.Request) {
	if !h.cfg.DemoMode {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	if err := demo.Apply(h.db, h.cfg); err != nil {
		log.Printf("demo reset failed: %v", err)
		writeError(w, http.StatusInternalServerError, "reset failed")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
