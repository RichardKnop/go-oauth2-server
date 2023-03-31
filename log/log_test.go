package log_test

import (
	"testing"

	"github.com/caputomarcos/go-oauth2-server/log"
)

func TestDefaultLogger(t *testing.T) {
	log.INFO.Print("should not panic")
	log.WARNING.Print("should not panic")
	log.ERROR.Print("should not panic")
	log.FATAL.Print("should not panic")
}
