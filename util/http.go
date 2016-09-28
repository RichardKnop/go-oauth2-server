package util

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// ParseBearerToken parses Bearer token from Authorization header
func ParseBearerToken(r *http.Request) ([]byte, error) {
	auth := r.Header.Get("Authorization")

	if !strings.HasPrefix(auth, "Bearer ") {
		return nil, errors.New("Bearer token not found")
	}

	bearerToken := strings.TrimPrefix(auth, "Bearer ")
	return []byte(bearerToken), nil
}

// GetCurrentURL returns the current request URL
func GetCurrentURL(r *http.Request) string {
	url := r.URL.Path
	qs := r.URL.Query().Encode()
	if qs != "" {
		url = fmt.Sprintf("%s?%s", url, qs)
	}
	return url
}
