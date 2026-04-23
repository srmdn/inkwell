package handler

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/srmdn/foliocms/internal/middleware"
)

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// ChangePassword verifies the current password and replaces it with a new one.
// Clears the session cookie so the user must re-authenticate.
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	var req changePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.CurrentPassword == "" || req.NewPassword == "" {
		writeError(w, http.StatusBadRequest, "current_password and new_password are required")
		return
	}
	if len(req.NewPassword) < 8 {
		writeError(w, http.StatusBadRequest, "new password must be at least 8 characters")
		return
	}

	var passwdHash string
	err := h.db.QueryRow(`SELECT passwd_hash FROM users WHERE id = ?`, userID).Scan(&passwdHash)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not fetch user")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwdHash), []byte(req.CurrentPassword)); err != nil {
		writeError(w, http.StatusUnauthorized, "current password is incorrect")
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	if _, err := h.db.Exec(`UPDATE users SET passwd_hash = ? WHERE id = ?`, string(newHash), userID); err != nil {
		writeError(w, http.StatusInternalServerError, "could not update password")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusNoContent)
}
