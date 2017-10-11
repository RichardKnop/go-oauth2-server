package response

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	thelog "github.com/RichardKnop/go-oauth2-server/log"
	"github.com/urfave/negroni"
)

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type Logger struct {
	// Logger inherits from log.Logger used to log messages with the Logger middleware
	*log.Logger
}

// NewURLLogger returns a new Logger instance
func NewURLLogger() *Logger {
	return &Logger{log.New(os.Stdout, "[negroni] ", 0)}
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	ip := r.RemoteAddr
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ip = xff
	}

	thelog.INFO.Printf("Started %s %s for %s", r.Method, r.URL.Path, ip)

	next(rw, r)

	res := rw.(negroni.ResponseWriter)

	msg := fmt.Sprintf("Finished %s %s : %v %s in %v", r.Method, r.URL.Path, res.Status(), http.StatusText(res.Status()), time.Since(start))

	switch {
	case res.Status() < 400:
		thelog.INFO.Print(msg)
	case res.Status() < 500:
		thelog.WARNING.Print(msg)
	default:
		thelog.ERROR.Print(msg)
	}
}
