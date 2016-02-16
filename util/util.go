package util

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/lib/pq"
)

// IntOrNull returns properly configured sql.NullInt64
func IntOrNull(n int64) sql.NullInt64 {
	return sql.NullInt64{Int64: n, Valid: true}
}

// PositiveIntOrNull returns properly configured sql.NullInt64 for a positive number
func PositiveIntOrNull(n int64) sql.NullInt64 {
	if n <= 0 {
		return sql.NullInt64{Int64: 0, Valid: false}
	}
	return sql.NullInt64{Int64: n, Valid: true}
}

// FloatOrNull returns properly configured sql.NullFloat64
func FloatOrNull(n float64) sql.NullFloat64 {
	return sql.NullFloat64{Float64: n, Valid: true}
}

// PositiveFloatOrNull returns properly configured sql.NullFloat64 for a positive number
func PositiveFloatOrNull(n float64) sql.NullFloat64 {
	if n <= 0 {
		return sql.NullFloat64{Float64: 0.0, Valid: false}
	}
	return sql.NullFloat64{Float64: n, Valid: true}
}

// StringOrNull returns properly configured sql.NullString
func StringOrNull(str string) sql.NullString {
	if str == "" {
		return sql.NullString{String: "", Valid: false}
	}
	return sql.NullString{String: str, Valid: true}
}

// TimeOrNull returns properly confiigured pq.TimeNull
func TimeOrNull(t *time.Time) pq.NullTime {
	if t == nil {
		return pq.NullTime{Time: time.Time{}, Valid: false}
	}
	return pq.NullTime{Time: *t, Valid: true}
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
