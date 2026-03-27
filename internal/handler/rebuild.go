package handler

import (
	"net/http"

	"github.com/srmdn/inkwell/internal/rebuild"
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
