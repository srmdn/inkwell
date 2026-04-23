package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/srmdn/foliocms/internal/config"
	"github.com/srmdn/foliocms/internal/db"
	"github.com/srmdn/foliocms/internal/handler"
	"github.com/srmdn/foliocms/internal/middleware"
)

func newAccountHandler(t *testing.T) (*handler.Handler, *db.DB) {
	t.Helper()
	database := newTestDB(t)
	cfg := &config.Config{
		ContentDir: t.TempDir(),
		JWTSecret:  "test-secret",
	}
	return handler.New(database, cfg), database
}

func insertUser(t *testing.T, database *db.DB, email, password string) int64 {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("bcrypt: %v", err)
	}
	res, err := database.Exec(
		`INSERT INTO users (email, passwd_hash) VALUES (?, ?)`, email, string(hash),
	)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}
	id, _ := res.LastInsertId()
	return id
}

func routerWithChangePassword(h *handler.Handler, userID int64) *chi.Mux {
	r := chi.NewRouter()
	r.Put("/api/admin/account/password", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), middleware.ContextKeyUserID, userID)
		h.ChangePassword(w, r.WithContext(ctx))
	})
	return r
}

func changePasswordReq(t *testing.T, r *chi.Mux, current, next string) *httptest.ResponseRecorder {
	t.Helper()
	b, _ := json.Marshal(map[string]string{
		"current_password": current,
		"new_password":     next,
	})
	req := httptest.NewRequest(http.MethodPut, "/api/admin/account/password", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestChangePasswordSuccess(t *testing.T) {
	h, database := newAccountHandler(t)
	userID := insertUser(t, database, "admin@example.com", "oldpass1")
	r := routerWithChangePassword(h, userID)

	w := changePasswordReq(t, r, "oldpass1", "newpass123")
	if w.Code != http.StatusNoContent {
		t.Fatalf("status: got %d, want %d\nbody: %s", w.Code, http.StatusNoContent, w.Body.String())
	}

	// Cookie should be cleared.
	cookies := w.Result().Cookies()
	for _, c := range cookies {
		if c.Name == "token" && c.MaxAge >= 0 {
			t.Errorf("expected token cookie to be cleared, got MaxAge=%d", c.MaxAge)
		}
	}

	// New password must work for login.
	var newHash string
	database.QueryRow(`SELECT passwd_hash FROM users WHERE id = ?`, userID).Scan(&newHash)
	if err := bcrypt.CompareHashAndPassword([]byte(newHash), []byte("newpass123")); err != nil {
		t.Errorf("new password hash does not match: %v", err)
	}
}

func TestChangePasswordWrongCurrent(t *testing.T) {
	h, database := newAccountHandler(t)
	userID := insertUser(t, database, "admin@example.com", "correctpass")
	r := routerWithChangePassword(h, userID)

	w := changePasswordReq(t, r, "wrongpass", "newpass123")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestChangePasswordTooShort(t *testing.T) {
	h, database := newAccountHandler(t)
	userID := insertUser(t, database, "admin@example.com", "oldpass1")
	r := routerWithChangePassword(h, userID)

	w := changePasswordReq(t, r, "oldpass1", "short")
	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestChangePasswordMissingFields(t *testing.T) {
	h, database := newAccountHandler(t)
	userID := insertUser(t, database, "admin@example.com", "oldpass1")
	r := routerWithChangePassword(h, userID)

	w := changePasswordReq(t, r, "", "newpass123")
	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}
