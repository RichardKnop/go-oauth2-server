package util

import (
	"github.com/unrolled/secure"
)

// NewSecure returns instance of secure.Secure to redirect http to https
func NewSecure(isDevelopment bool) *secure.Secure {
	return secure.New(secure.Options{
		SSLRedirect:     true,
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		IsDevelopment:   isDevelopment,
	})
}
