package handler

import (
	"encoding/json"
	"net/http"

	"github.com/srmdn/foliocms/internal/config"
	"github.com/srmdn/foliocms/internal/db"
	"github.com/srmdn/foliocms/internal/rebuild"
)

// Handler holds shared dependencies for all HTTP handlers.
type Handler struct {
	db        *db.DB
	cfg       *config.Config
	rebuilder *rebuild.Rebuilder
}

func New(database *db.DB, cfg *config.Config) *Handler {
	return &Handler{db: database, cfg: cfg}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
