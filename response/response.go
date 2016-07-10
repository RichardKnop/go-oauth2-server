package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const realm = "recall"

// WriteJSON writes JSON response
func WriteJSON(w http.ResponseWriter, v interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

// NoContent writes a 204 no content response
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Error produces a JSON error response with the following structure:
// {"error":"some error message"}
func Error(w http.ResponseWriter, err string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": err})
}

// UnauthorizedError has to contain WWW-Authenticate header
// See http://self-issued.info/docs/draft-ietf-oauth-v2-bearer.html#rfc.section.3
func UnauthorizedError(w http.ResponseWriter, err string) {
	// TODO - include error if the request contained an access token
	w.Header().Set("WWW-Authenticate", fmt.Sprintf("Bearer realm=%s", realm))
	Error(w, err, http.StatusUnauthorized)
}
