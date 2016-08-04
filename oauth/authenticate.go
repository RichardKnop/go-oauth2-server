package oauth

import (
	"errors"
	"time"
)

var (
	// ErrAccessTokenNotFound ...
	ErrAccessTokenNotFound = errors.New("Access token not found")
	// ErrAccessTokenExpired ...
	ErrAccessTokenExpired = errors.New("Access token expired")
)

// Authenticate checks the access token is valid
func (s *Service) Authenticate(token string) (*AccessToken, error) {
	// Fetch the access token from the database
	accessToken := new(AccessToken)
	notFound := s.db.Where("token = ?", token).First(accessToken).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrAccessTokenNotFound
	}

	// Check the access token hasn't expired
	if time.Now().After(accessToken.ExpiresAt) {
		return nil, ErrAccessTokenExpired
	}

	// Extend refresh token expiration database
	query := s.db.Model(new(RefreshToken)).Where("client_id = ?", accessToken.ClientID.Int64)
	if accessToken.UserID.Valid {
		query = query.Where("user_id = ?", accessToken.UserID.Int64)
	} else {
		query = query.Where("user_id IS NULL")
	}
	increasedExpiresAt := time.Now().Add(
		time.Duration(s.cnf.Oauth.RefreshTokenLifetime) * time.Second,
	)
	if err := query.UpdateColumn("expires_at", increasedExpiresAt).Error; err != nil {
		return nil, err
	}

	return accessToken, nil
}
