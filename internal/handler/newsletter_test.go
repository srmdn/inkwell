package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/srmdn/foliocms/internal/config"
	"github.com/srmdn/foliocms/internal/handler"
)

type recordingSender struct {
	calls []sentEmail
}

type sentEmail struct {
	to, subject, body string
}

func (r *recordingSender) Send(to, subject, body string) error {
	r.calls = append(r.calls, sentEmail{to, subject, body})
	return nil
}

func newsletterRouter(h *handler.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/api/subscribe", h.Subscribe)
	r.Post("/api/admin/newsletter/send", h.SendNewsletter)
	return r
}

func newNewsletterHandler(t *testing.T, dir string, siteURL string) (*handler.Handler, *recordingSender) {
	t.Helper()
	database := newTestDB(t)
	cfg := &config.Config{
		ContentDir: dir,
		JWTSecret:  "test-secret",
		SMTPHost:   "smtp.example.com",
		SiteURL:    siteURL,
	}
	h := handler.New(database, cfg)
	sender := &recordingSender{}
	h.SetMailer(sender)
	return h, sender
}

func subscribeEmail(t *testing.T, r *chi.Mux, email string) {
	t.Helper()
	b, _ := json.Marshal(map[string]string{"email": email})
	req := httptest.NewRequest(http.MethodPost, "/api/subscribe", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("subscribe %s: got %d", email, w.Code)
	}
}

func TestSendNewsletterInjectsUnsubscribeLink(t *testing.T) {
	dir := t.TempDir()
	h, sender := newNewsletterHandler(t, dir, "https://example.com")
	r := newsletterRouter(h)

	subscribeEmail(t, r, "alice@example.com")
	subscribeEmail(t, r, "bob@example.com")

	payload, _ := json.Marshal(map[string]string{
		"subject": "Hello",
		"body":    "Newsletter body.",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/admin/newsletter/send", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("send: got %d, body: %s", w.Code, w.Body.String())
	}
	if len(sender.calls) != 2 {
		t.Fatalf("expected 2 sends, got %d", len(sender.calls))
	}
	for _, call := range sender.calls {
		if !strings.Contains(call.body, "https://example.com/api/unsubscribe?token=") {
			t.Errorf("email to %s missing unsubscribe link; body: %s", call.to, call.body)
		}
		if !strings.Contains(call.body, "Newsletter body.") {
			t.Errorf("email to %s missing original body", call.to)
		}
	}
}

func TestSendNewsletterNoSubscribers(t *testing.T) {
	dir := t.TempDir()
	h, sender := newNewsletterHandler(t, dir, "https://example.com")
	r := newsletterRouter(h)

	payload, _ := json.Marshal(map[string]string{"subject": "Hi", "body": "Body."})
	req := httptest.NewRequest(http.MethodPost, "/api/admin/newsletter/send", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("got %d", w.Code)
	}
	if len(sender.calls) != 0 {
		t.Errorf("expected 0 sends, got %d", len(sender.calls))
	}
}

func TestSendNewsletterSMTPNotConfigured(t *testing.T) {
	dir := t.TempDir()
	h := newTestHandler(t, dir) // SMTPHost is empty
	r := newsletterRouter(h)

	payload, _ := json.Marshal(map[string]string{"subject": "Hi", "body": "Body."})
	req := httptest.NewRequest(http.MethodPost, "/api/admin/newsletter/send", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", w.Code)
	}
}
