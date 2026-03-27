package middleware

import (
	"net/http"

	"github.com/srmdn/foliocms/internal/auth"
)

// VerifyCSRF validates the X-CSRF-Token header on mutating requests.
// It must be used after Authenticate, since it reads claims from context.
//
// The expected token is HMAC-SHA256(jti, jwtSecret) — derived from the JWT,
// no extra storage required. The client receives the CSRF token in the login
// response and must include it on every POST/PUT/PATCH/DELETE request.
func VerifyCSRF(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet, http.MethodHead, http.MethodOptions:
				next.ServeHTTP(w, r)
				return
			}

			claims := ClaimsFromContext(r.Context())
			if claims == nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			provided := r.Header.Get("X-CSRF-Token")
			expected := auth.CSRFTokenFromClaims(claims, jwtSecret)

			if provided == "" || provided != expected {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
