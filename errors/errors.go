package errors

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
)

// Error produces an error response in JSON with the following structure, '{"error":"My error message"}'
// The standard plain text net/http Error helper can still be called like this:
// http.Error(w, "error message", code)
func Error(w rest.ResponseWriter, error string, code int) {
	w.WriteHeader(code)
	err := w.WriteJson(map[string]string{"error": error})
	if err != nil {
		panic(err)
	}
}

// UnauthorizedError has to contain WWW-Authenticate header
// See http://self-issued.info/docs/draft-ietf-oauth-v2-bearer.html#rfc.section.3
func UnauthorizedError(w rest.ResponseWriter, err string) {
	// TODO - include error if the request contained an access token
	w.Header().Set("WWW-Authenticate", "Bearer realm=areatech_api")
	Error(w, err, http.StatusUnauthorized)
}
