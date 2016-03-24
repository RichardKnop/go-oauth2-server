package oauth

import (
	"errors"
	"time"

	"github.com/RichardKnop/go-oauth2-server/util"
)

var (
	// ErrRefreshTokenNotFound ...
	ErrRefreshTokenNotFound = errors.New("Refresh token not found")
	// ErrRefreshTokenExpired ...
	ErrRefreshTokenExpired = errors.New("Refresh token expired")
)

// GetOrCreateRefreshToken retrieves an existing refresh token, if expired,
// the token gets deleted and new refresh token is created
func (s *Service) GetOrCreateRefreshToken(client *Client, user *User, expiresIn int, scope string) (*RefreshToken, error) {
	// Try to fetch an existing refresh token first
	refreshToken := new(RefreshToken)
	clientID := util.PositiveIntOrNull(int64(client.ID))
	userID := util.PositiveIntOrNull(0) // user ID can be NULL
	if user != nil {
		userID = util.PositiveIntOrNull(int64(user.ID))
	}
	found := !s.db.Where(RefreshToken{
		ClientID: clientID,
		UserID:   userID,
	}).Preload("Client").Preload("User").First(refreshToken).RecordNotFound()

	// Check if the token is expired, if found
	var expired bool
	if found {
		expired = time.Now().After(refreshToken.ExpiresAt)
	}

	// If the refresh token has expired, delete it
	if expired {
		s.db.Unscoped().Delete(refreshToken)
	}

	// Create a new refresh token if it expired or was not found
	if expired || !found {
		refreshToken = NewRefreshToken(client, user, expiresIn, scope)
		if err := s.db.Create(refreshToken).Error; err != nil {
			return nil, err
		}
	}

	return refreshToken, nil
}

// GetValidRefreshToken returns a valid non expired refresh token
func (s *Service) GetValidRefreshToken(token string, client *Client) (*RefreshToken, error) {
	// Fetch the refresh token from the database
	refreshToken := new(RefreshToken)
	notFound := s.db.Where(RefreshToken{
		ClientID: util.PositiveIntOrNull(int64(client.ID)),
	}).Where("token = ?", token).Preload("Client").Preload("User").
		First(refreshToken).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrRefreshTokenNotFound
	}

	// Check the refresh token hasn't expired
	if time.Now().After(refreshToken.ExpiresAt) {
		return nil, ErrRefreshTokenExpired
	}

	return refreshToken, nil
}
