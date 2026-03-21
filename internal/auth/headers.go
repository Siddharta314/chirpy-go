package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHedear := headers.Get("Authorization")
	if authHedear == "" {
		return "", errors.New("no authorization header")
	}
	authHedear = strings.TrimSpace(authHedear)
	if !strings.HasPrefix(authHedear, "Bearer ") {
		return "", errors.New("invalid authorization header")
	}
	return strings.TrimPrefix(authHedear, "Bearer "), nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authHedear := headers.Get("Authorization")
	if authHedear == "" {
		return "", errors.New("no api key header")
	}
	authHedear = strings.TrimSpace(authHedear)
	if !strings.HasPrefix(authHedear, "ApiKey ") {
		return "", errors.New("invalid api key header")
	}
	return strings.TrimPrefix(authHedear, "ApiKey "), nil
}
