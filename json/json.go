package json

import (
	"encoding/json"
	"net/http"
)

// WriteJSON writes JSON response
func WriteJSON(w http.ResponseWriter, v interface{}, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(v)
}

// Error produces a JSON error response with the following structure:
// {"error":"some error message"}
func Error(w http.ResponseWriter, err string, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"error": err})
}

// UnauthorizedError has to contain WWW-Authenticate header
// See http://self-issued.info/docs/draft-ietf-oauth-v2-bearer.html#rfc.section.3
func UnauthorizedError(w http.ResponseWriter, err string) {
	// TODO - include error if the request contained an access token
	w.Header().Set("WWW-Authenticate", "Bearer realm=areatech_api")
	Error(w, err, http.StatusUnauthorized)
}
