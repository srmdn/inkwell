package auth_test

import (
	"testing"

	"github.com/srmdn/foliocms/internal/auth"
)

const testSecret = "test-secret-not-for-production"

func TestGenerateAndParseToken(t *testing.T) {
	token, err := auth.GenerateToken(42, testSecret)
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := auth.ParseToken(token, testSecret)
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}

	if claims.UserID != 42 {
		t.Errorf("UserID: got %d, want 42", claims.UserID)
	}
	if claims.ID == "" {
		t.Error("expected non-empty jti")
	}
}

func TestParseTokenWrongSecret(t *testing.T) {
	token, err := auth.GenerateToken(1, testSecret)
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	_, err = auth.ParseToken(token, "wrong-secret")
	if err == nil {
		t.Error("expected error when parsing token with wrong secret")
	}
}

func TestParseTokenInvalid(t *testing.T) {
	_, err := auth.ParseToken("not.a.token", testSecret)
	if err == nil {
		t.Error("expected error for invalid token string")
	}
}

func TestCSRFTokenDeterministic(t *testing.T) {
	token, _ := auth.GenerateToken(1, testSecret)
	claims, _ := auth.ParseToken(token, testSecret)

	csrf1 := auth.CSRFTokenFromClaims(claims, testSecret)
	csrf2 := auth.CSRFTokenFromClaims(claims, testSecret)

	if csrf1 == "" {
		t.Fatal("expected non-empty CSRF token")
	}
	if csrf1 != csrf2 {
		t.Error("CSRF token should be deterministic for the same claims and secret")
	}
}

func TestCSRFTokenDiffersAcrossJWTs(t *testing.T) {
	token1, _ := auth.GenerateToken(1, testSecret)
	token2, _ := auth.GenerateToken(1, testSecret)

	claims1, _ := auth.ParseToken(token1, testSecret)
	claims2, _ := auth.ParseToken(token2, testSecret)

	csrf1 := auth.CSRFTokenFromClaims(claims1, testSecret)
	csrf2 := auth.CSRFTokenFromClaims(claims2, testSecret)

	if csrf1 == csrf2 {
		t.Error("CSRF tokens for different JWTs (different jti) should differ")
	}
}

func TestCSRFTokenDiffersWithDifferentSecret(t *testing.T) {
	token, _ := auth.GenerateToken(1, testSecret)
	claims, _ := auth.ParseToken(token, testSecret)

	csrf1 := auth.CSRFTokenFromClaims(claims, testSecret)
	csrf2 := auth.CSRFTokenFromClaims(claims, "different-secret")

	if csrf1 == csrf2 {
		t.Error("CSRF tokens should differ when secret differs")
	}
}
