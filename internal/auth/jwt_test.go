package auth

import (
	"testing"
	"time"
)

func TestGenerateAndValidateToken(t *testing.T) {
	secret := "test-secret-key"
	token, err := GenerateToken("user-123", "flakerimi", false, secret, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatal(err)
	}

	if claims.UserID != "user-123" {
		t.Errorf("UserID = %q, want %q", claims.UserID, "user-123")
	}
	if claims.Username != "flakerimi" {
		t.Errorf("Username = %q, want %q", claims.Username, "flakerimi")
	}
	if claims.IsAdmin {
		t.Error("IsAdmin = true, want false")
	}
}

func TestValidateTokenWrongSecret(t *testing.T) {
	token, _ := GenerateToken("user-123", "flakerimi", false, "secret-a", time.Hour)
	_, err := ValidateToken(token, "secret-b")
	if err == nil {
		t.Error("expected error for wrong secret")
	}
}

func TestValidateTokenExpired(t *testing.T) {
	token, _ := GenerateToken("user-123", "flakerimi", false, "secret", -time.Hour)
	_, err := ValidateToken(token, "secret")
	if err == nil {
		t.Error("expected error for expired token")
	}
}
