package oauth

import (
	"net/http"
)

var (
	errStatusCodeMap = map[error]int{
		ErrAuthorizationCodeNotFound:     http.StatusNotFound,
		ErrAuthorizationCodeExpired:      http.StatusBadRequest,
		ErrInvalidRedirectURI:            http.StatusBadRequest,
		ErrInvalidScope:                  http.StatusBadRequest,
		ErrInvalidUsernameOrPassword:     http.StatusBadRequest,
		ErrRefreshTokenNotFound:          http.StatusNotFound,
		ErrRefreshTokenExpired:           http.StatusBadRequest,
		ErrRequestedScopeCannotBeGreater: http.StatusBadRequest,
		ErrTokenMissing:                  http.StatusNotFound,
		ErrTokenHintInvalid:              http.StatusBadRequest,
		ErrAccessTokenNotFound:           http.StatusNotFound,
		ErrRefreshTokenNotFound:          http.StatusNotFound,
		ErrTokenMissing:                  http.StatusBadRequest,
		ErrTokenHintInvalid:              http.StatusBadRequest,
		ErrInvalidUsernameOrPassword:     http.StatusUnauthorized,
	}
)

func getErrStatusCode(err error) int {
	code, ok := errStatusCodeMap[err]
	if ok {
		return code
	}

	return http.StatusInternalServerError
}
