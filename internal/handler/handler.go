package handler

import (
	"encoding/json"
	"net/http"

	"github.com/srmdn/foliocms/internal/config"
	"github.com/srmdn/foliocms/internal/db"
	"github.com/srmdn/foliocms/internal/mailer"
	"github.com/srmdn/foliocms/internal/media"
	"github.com/srmdn/foliocms/internal/rebuild"
)

// Handler holds shared dependencies for all HTTP handlers.
type Handler struct {
	db          *db.DB
	cfg         *config.Config
	rebuilder   *rebuild.Rebuilder
	mailer      *mailer.Mailer
	mediaDriver media.MediaDriver
}

func New(database *db.DB, cfg *config.Config) *Handler {
	m := mailer.New(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUsername, cfg.SMTPPassword, cfg.SMTPFrom)
	return &Handler{db: database, cfg: cfg, mailer: m}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
