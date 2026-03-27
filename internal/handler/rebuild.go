package handler

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/srmdn/foliocms/internal/rebuild"
)

func (h *Handler) TriggerRebuild(w http.ResponseWriter, r *http.Request) {
	started := h.rebuilder.Trigger()
	if !started {
		writeJSON(w, http.StatusConflict, map[string]string{
			"error": "rebuild already in progress",
		})
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) RebuildStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.rebuilder.GetStatus())
}

// SetRebuilder wires the Rebuilder into the Handler after construction.
// Called from main before routes are registered.
func (h *Handler) SetRebuilder(rb *rebuild.Rebuilder) {
	h.rebuilder = rb
}

// WebhookRebuild triggers a rebuild via a shared secret token.
// Disabled (404) if WEBHOOK_SECRET is not configured.
// Accepts the secret via X-Webhook-Secret header or Authorization: Bearer <token>.
func (h *Handler) WebhookRebuild(w http.ResponseWriter, r *http.Request) {
	if h.cfg.WebhookSecret == "" {
		http.NotFound(w, r)
		return
	}

	token := r.Header.Get("X-Webhook-Secret")
	if token == "" {
		auth := r.Header.Get("Authorization")
		token = strings.TrimPrefix(auth, "Bearer ")
	}

	if subtle.ConstantTimeCompare([]byte(token), []byte(h.cfg.WebhookSecret)) != 1 {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid secret"})
		return
	}

	started := h.rebuilder.Trigger()
	if !started {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "rebuild already in progress"})
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
