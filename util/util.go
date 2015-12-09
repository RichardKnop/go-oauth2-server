package util

import (
	"database/sql"
	"strings"
)

// IntOrNull returns properly confiigured sql.NullInt64
func IntOrNull(n uint) sql.NullInt64 {
	if n < 1 {
		return sql.NullInt64{Int64: 0, Valid: false}
	}

	return sql.NullInt64{Int64: int64(n), Valid: true}
}

// StringOrNull returns properly confiigured sql.NullString
func StringOrNull(str string) sql.NullString {
	if str == "" {
		return sql.NullString{String: "", Valid: false}
	}

	return sql.NullString{String: str, Valid: true}
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
