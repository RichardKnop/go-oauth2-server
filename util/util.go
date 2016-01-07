package util

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/lib/pq"
)

// IntOrNull returns properly confiigured sql.NullInt64
func IntOrNull(n int64) sql.NullInt64 {
	if n < 1 {
		return sql.NullInt64{Int64: 0, Valid: false}
	}

	return sql.NullInt64{Int64: n, Valid: true}
}

// StringOrNull returns properly confiigured sql.NullString
func StringOrNull(str string) sql.NullString {
	if str == "" {
		return sql.NullString{String: "", Valid: false}
	}

	return sql.NullString{String: str, Valid: true}
}

// TimeOrNull returns properly confiigured pq.TimeNull
func TimeOrNull(t interface{}) pq.NullTime {
	nullTime := new(pq.NullTime)
	nullTime.Scan(t)
	return *nullTime
}

// StringInSlice is a function similar to "x in y" Python construct
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// SpaceDelimitedStringNotGreater returns true if the first string
// is the same as the second string or does not contain any substring
// not contained in the second string (when split by space)
func SpaceDelimitedStringNotGreater(first, second string) bool {
	// Empty string is never greater
	if first == "" {
		return true
	}

	// Split the second string by space
	secondParts := strings.Split(second, " ")

	// Iterate over space delimited parts of the first string
	for _, firstPart := range strings.Split(first, " ") {
		// If the substring is not part of the second string, return false
		if !StringInSlice(firstPart, secondParts) {
			return false
		}
	}

	// The first string is the same or more restrictive
	// than the second string, return true
	return true
}

// ParseBearerToken parses Bearer token from Authorization header
func ParseBearerToken(r *http.Request) ([]byte, error) {
	auth := r.Header.Get("Authorization")

	if !strings.HasPrefix(auth, "Bearer ") {
		return nil, errors.New("Bearer token not found")
	}

	bearerToken := strings.TrimPrefix(auth, "Bearer ")
	return []byte(bearerToken), nil
}
