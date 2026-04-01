package handler

import (
	"encoding/json"
	"net/http"
)

// allowedSettingKeys is the whitelist of keys that can be read and updated.
var allowedSettingKeys = map[string]bool{
	"site_name":        true,
	"site_description": true,
	"social_github":    true,
	"social_twitter":   true,
	"social_linkedin":  true,
}

// GetPublicSettings returns all site settings (public endpoint, for themes).
func (h *Handler) GetPublicSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.loadSettings()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load settings")
		return
	}
	writeJSON(w, http.StatusOK, settings)
}

// GetAdminSettings returns all site settings (protected).
func (h *Handler) GetAdminSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.loadSettings()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load settings")
		return
	}
	writeJSON(w, http.StatusOK, settings)
}

// UpdateSettings updates one or more site settings (protected + CSRF).
func (h *Handler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	var updates map[string]string
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	for key, value := range updates {
		if !allowedSettingKeys[key] {
			writeError(w, http.StatusBadRequest, "unknown setting key: "+key)
			return
		}
		if _, err := h.db.Exec(`UPDATE settings SET value = ? WHERE key = ?`, value, key); err != nil {
			writeError(w, http.StatusInternalServerError, "could not update setting")
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) loadSettings() (map[string]string, error) {
	rows, err := h.db.Query(`SELECT key, value FROM settings`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		result[key] = value
	}
	return result, rows.Err()
}
