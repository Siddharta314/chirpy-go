package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	tokenSecret := "+3sYRsmi8xCFdKPVpAYcERXbLsVWGt/7o+IvJLYJ3Ik="
	userID := uuid.New()
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}

	parsedID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("failed to validate JWT: %v", err)
	}

	if parsedID != userID {
		t.Errorf("expected ID %v, got %v", userID, parsedID)
	}

	_, err = ValidateJWT(token, "wrong-secret")
	if err == nil {
		t.Error("expected error with wrong secret, got nil")
	}

	expiredToken, err := MakeJWT(userID, tokenSecret, -time.Hour)
	if err != nil {
		t.Fatalf("failed to create expired token: %v", err)
	}

	_, err = ValidateJWT(expiredToken, tokenSecret)
	if err == nil {
		t.Error("expected error with expired token, got nil")
	}
}