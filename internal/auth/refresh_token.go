package auth

import (
	"crypto/rand"
	"encoding/hex"
)


func MakeRefreshToken() string{
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}
