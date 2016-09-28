package util

import (
	"time"
)

// https://golang.org/pkg/time/#pkg-constants
const (
	// RFC3339Mili is a modification of RFC3339Nano to only include ms (3 decimal places)
	RFC3339Mili = "2006-01-02T15:04:05.999Z07:00"
	// DateFormat used for things such as date of birth (when time does not matter)
	DateFormat = "2006-01-02"
)

// FormatTime formats a time object to RFC3339 with ms precision
func FormatTime(timestamp *time.Time) string {
	if timestamp == nil {
		return ""
	}
	return timestamp.UTC().Format(RFC3339Mili)
}

// ParseTimestamp parses a string representation of a timestamp in RFC3339
// format and returns a time.Time instance
func ParseTimestamp(timestamp string) (*time.Time, error) {
	// RFC3339 = "2006-01-02T15:04:05Z07:00"
	if timestamp == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// FormatDate formats a time object to a date only format, e.g. 2006-01-02
func FormatDate(timestamp *time.Time) string {
	if timestamp == nil {
		return ""
	}
	return timestamp.UTC().Format(DateFormat)
}

// ParseDate parses a string representation of a date format
// and returns a time.Time instance
func ParseDate(timestamp string) (*time.Time, error) {
	if timestamp == "" {
		return nil, nil
	}
	t, err := time.Parse(DateFormat, timestamp)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
