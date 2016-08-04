package response

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestResponseForError tests a response w to see if it returned an error msg with http code
func TestResponseForError(t *testing.T, w *httptest.ResponseRecorder, msg string, code int) {
	assert.NotNil(t, w)
	assert.Equal(
		t,
		code,
		w.Code,
		fmt.Sprintf("Expected a %d response but got %d", code, w.Code),
	)
	TestResponseBody(t, w, getErrorJSON(msg))
}

// TestEmptyResponse tests an empty 204 response
func TestEmptyResponse(t *testing.T, w *httptest.ResponseRecorder) {
	assert.Equal(t, 204, w.Code)
	TestResponseBody(t, w, "")
}

// TestResponseObject tests response body is equal to expected object in JSON form
func TestResponseObject(t *testing.T, w *httptest.ResponseRecorder, expected interface{}, code int) {
	assert.Equal(
		t,
		code,
		w.Code,
		fmt.Sprintf("Expected a %d response but got %d", code, w.Code),
	)
	jsonBytes, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.Equal(
		t,
		string(jsonBytes),
		strings.TrimRight(w.Body.String(), "\n"),
		"Should have returned correct body text",
	)
}

// TestResponseBody tests response body is equal to expected string
func TestResponseBody(t *testing.T, w *httptest.ResponseRecorder, expected string) {
	assert.Equal(
		t,
		expected,
		strings.TrimRight(w.Body.String(), "\n"),
		"Should have returned correct body text",
	)
}

func getErrorJSON(msg string) string {
	return fmt.Sprintf("{\"error\":\"%s\"}", msg)
}
