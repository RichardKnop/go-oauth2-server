package oauth

import (
	"errors"
	"time"
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
	query := RefreshTokenPreload(s.db).Where("client_id = ?", client.ID)
	if user != nil && user.ID > 0 {
		query = query.Where("user_id = ?", user.ID)
	} else {
		query = query.Where("user_id IS NULL")
	}
	found := !query.First(refreshToken).RecordNotFound()

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
		refreshToken.Client = client
		refreshToken.User = user
	}

	return refreshToken, nil
}

// GetValidRefreshToken returns a valid non expired refresh token
func (s *Service) GetValidRefreshToken(token string, client *Client) (*RefreshToken, error) {
	// Fetch the refresh token from the database
	refreshToken := new(RefreshToken)
	notFound := RefreshTokenPreload(s.db).Where("client_id = ?", client.ID).
		Where("token = ?", token).First(refreshToken).RecordNotFound()

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
