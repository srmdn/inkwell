package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/srmdn/foliocms/internal/auth"
)

// Authenticate validates the JWT from either the "token" cookie or the
// Authorization: Bearer header. On success it injects the claims into context.
// On failure it returns 401 and stops the chain.
func Authenticate(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := tokenFromRequest(r)
			if tokenStr == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := auth.ParseToken(tokenStr, jwtSecret)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
			ctx = context.WithValue(ctx, ContextKeyClaims, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserIDFromContext retrieves the authenticated user's ID from context.
// Returns 0 if not set.
func UserIDFromContext(ctx context.Context) int64 {
	id, _ := ctx.Value(ContextKeyUserID).(int64)
	return id
}

// ClaimsFromContext retrieves the JWT claims from context.
func ClaimsFromContext(ctx context.Context) *auth.Claims {
	claims, _ := ctx.Value(ContextKeyClaims).(*auth.Claims)
	return claims
}

func tokenFromRequest(r *http.Request) string {
	// Prefer cookie (browser dashboard flow)
	if cookie, err := r.Cookie("token"); err == nil && cookie.Value != "" {
		return cookie.Value
	}
	// Fall back to Authorization header (API / CLI clients)
	if header := r.Header.Get("Authorization"); strings.HasPrefix(header, "Bearer ") {
		return strings.TrimPrefix(header, "Bearer ")
	}
	return ""
}
