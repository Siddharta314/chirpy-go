package auth

import (
	"testing"
)

func TestPasswordHashing(t *testing.T) {
	password := "my-secure-password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if hash == password {
		t.Fatal("hash should not be equal to plain text password")
	}

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("failed to check password hash: %v", err)
	}
	if !match {
		t.Error("expected password to match its hash")
	}

	match, err = CheckPasswordHash("wrong-password", hash)
	if err != nil {
		t.Fatalf("failed to check password hash: %v", err)
	}
	if match {
		t.Error("expected wrong password to not match the hash")
	}
}