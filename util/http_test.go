package util_test

import (
	"net/http"
	"testing"

	"github.com/RichardKnop/go-oauth2-server/util"
	"github.com/stretchr/testify/assert"
)

func TestParseBearerTokenNotFound(t *testing.T) {
	r, err := http.NewRequest("GET", "http://1.2.3.4/something", nil)
	assert.NoError(t, err, "Request setup should not get an error")
	r.Header.Add("Authorization", "bogus bogus")

	token, err := util.ParseBearerToken(r)

	// Token should be nil
	assert.Nil(t, token)

	// Correct error should be returned
	if assert.NotNil(t, err) {
		assert.Equal(t, "Bearer token not found", err.Error())
	}
}

func TestParseBearerToken(t *testing.T) {
	r, err := http.NewRequest("GET", "http://1.2.3.4/something", nil)
	assert.NoError(t, err, "Request setup should not get an error")
	r.Header.Add("Authorization", "Bearer test_token")

	token, err := util.ParseBearerToken(r)

	// Error should be nil
	assert.Nil(t, err)

	// Correct token should be returned
	if assert.NotNil(t, token) {
		assert.Equal(t, []byte("test_token"), token)
	}
}
