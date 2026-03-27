package middleware

type contextKey string

const (
	ContextKeyUserID contextKey = "userID"
	ContextKeyClaims contextKey = "claims"
)
