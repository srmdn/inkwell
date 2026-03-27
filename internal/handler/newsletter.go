package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/srmdn/foliocms/internal/model"
)

// Subscribe adds a new subscriber. POST /api/subscribe
func (h *Handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		writeError(w, http.StatusBadRequest, "email is required")
		return
	}

	token := uuid.New().String()
	_, err := h.db.Exec(
		`INSERT INTO subscribers (email, token) VALUES (?, ?)`,
		req.Email, token,
	)
	if err != nil {
		// SQLite unique constraint violation
		if isUniqueErr(err) {
			writeError(w, http.StatusConflict, "already subscribed")
			return
		}
		writeError(w, http.StatusInternalServerError, "could not subscribe")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Unsubscribe removes a subscriber by token. GET /api/unsubscribe?token=xxx
func (h *Handler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		writeError(w, http.StatusBadRequest, "token is required")
		return
	}

	res, err := h.db.Exec(`DELETE FROM subscribers WHERE token = ?`, token)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not unsubscribe")
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		writeError(w, http.StatusNotFound, "token not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListSubscribers returns all subscribers. GET /api/admin/subscribers
func (h *Handler) ListSubscribers(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(
		`SELECT id, email, token, subscribed_at FROM subscribers ORDER BY subscribed_at DESC`,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not list subscribers")
		return
	}
	defer rows.Close()

	var subs []model.Subscriber
	for rows.Next() {
		var s model.Subscriber
		if err := rows.Scan(&s.ID, &s.Email, &s.Token, &s.SubscribedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "scan error")
			return
		}
		subs = append(subs, s)
	}
	if subs == nil {
		subs = []model.Subscriber{}
	}
	writeJSON(w, http.StatusOK, subs)
}

// DeleteSubscriber removes a subscriber by ID. DELETE /api/admin/subscribers/{id}
func (h *Handler) DeleteSubscriber(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	res, err := h.db.Exec(`DELETE FROM subscribers WHERE id = ?`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not delete subscriber")
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		writeError(w, http.StatusNotFound, "subscriber not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// SendNewsletter sends an email to all subscribers. POST /api/admin/newsletter/send
func (h *Handler) SendNewsletter(w http.ResponseWriter, r *http.Request) {
	if h.cfg.SMTPHost == "" {
		writeError(w, http.StatusServiceUnavailable, "SMTP not configured")
		return
	}

	var req struct {
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Subject == "" || req.Body == "" {
		writeError(w, http.StatusBadRequest, "subject and body are required")
		return
	}

	rows, err := h.db.Query(`SELECT email FROM subscribers`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not fetch subscribers")
		return
	}
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			writeError(w, http.StatusInternalServerError, "scan error")
			return
		}
		emails = append(emails, email)
	}

	if len(emails) == 0 {
		writeJSON(w, http.StatusOK, map[string]any{"sent": 0, "message": "no subscribers"})
		return
	}

	if err := h.mailer.SendMany(emails, req.Subject, req.Body); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"sent": len(emails)})
}

func isUniqueErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed")
}
