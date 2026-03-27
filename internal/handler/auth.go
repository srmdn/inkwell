package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/srmdn/foliocms/internal/auth"
	"github.com/srmdn/foliocms/internal/middleware"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	CSRFToken string `json:"csrf_token"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	var userID int64
	var passwdHash string
	err := h.db.QueryRow(
		`SELECT id, passwd_hash FROM users WHERE email = ?`, req.Email,
	).Scan(&userID, &passwdHash)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwdHash), []byte(req.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := auth.GenerateToken(userID, h.cfg.JWTSecret)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not generate token")
		return
	}

	// Parse back to get claims (needed to derive CSRF token)
	claims, err := auth.ParseToken(token, h.cfg.JWTSecret)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not read token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	writeJSON(w, http.StatusOK, loginResponse{
		CSRFToken: auth.CSRFTokenFromClaims(claims, h.cfg.JWTSecret),
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
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

// GetCSRFToken returns the CSRF token for the currently authenticated session.
// Clients that lose the CSRF token (e.g. page refresh) can call this to
// retrieve it without logging in again.
func (h *Handler) GetCSRFToken(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"csrf_token": auth.CSRFTokenFromClaims(claims, h.cfg.JWTSecret),
	})
}
