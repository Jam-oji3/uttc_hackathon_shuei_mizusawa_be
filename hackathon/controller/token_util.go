package controller

import (
	"errors"
	"net/http"
	"strings"
)

var (
	ErrNoAuthHeader      = errors.New("authorization header is required")
	ErrInvalidAuthHeader = errors.New("invalid authorization header format")
)

func ExtractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeader
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", ErrInvalidAuthHeader
	}
	return parts[1], nil
}
