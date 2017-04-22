package oauth

import (
	"errors"
	"time"

	"github.com/RichardKnop/go-oauth2-server/models"
	"github.com/RichardKnop/go-oauth2-server/session"
	"github.com/jinzhu/gorm"
)

var (
	// ErrAccessTokenNotFound ...
	ErrAccessTokenNotFound = errors.New("Access token not found")
	// ErrAccessTokenExpired ...
	ErrAccessTokenExpired = errors.New("Access token expired")
)

// Authenticate checks the access token is valid
func (s *Service) Authenticate(token string) (*models.OauthAccessToken, error) {
	// Fetch the access token from the database
	accessToken := new(models.OauthAccessToken)
	notFound := s.db.Where("token = ?", token).First(accessToken).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrAccessTokenNotFound
	}

	// Check the access token hasn't expired
	if time.Now().UTC().After(accessToken.ExpiresAt) {
		return nil, ErrAccessTokenExpired
	}

	// Extend refresh token expiration database
	query := s.db.Model(new(models.OauthRefreshToken)).Where("client_id = ?", accessToken.ClientID.String)
	if accessToken.UserID.Valid {
		query = query.Where("user_id = ?", accessToken.UserID.String)
	} else {
		query = query.Where("user_id IS NULL")
	}
	increasedExpiresAt := gorm.NowFunc().Add(
		time.Duration(s.cnf.Oauth.RefreshTokenLifetime) * time.Second,
	)
	if err := query.UpdateColumn("expires_at", increasedExpiresAt).Error; err != nil {
		return nil, err
	}

	return accessToken, nil
}

// ClearUserTokens deletes the user's access and refresh tokens associated with this client id
func (s *Service) ClearUserTokens(userSession *session.UserSession) {
	// Clear all refresh tokens with user_id and client_id
	refreshToken := new(models.OauthRefreshToken)
	found := !models.OauthRefreshTokenPreload(s.db).Where("token = ?", userSession.RefreshToken).First(refreshToken).RecordNotFound()
	if found {
		s.db.Unscoped().Where("client_id = ? AND user_id = ?", refreshToken.ClientID, refreshToken.UserID).Delete(models.OauthRefreshToken{})
	}

	// Clear all access tokens with user_id and client_id
	accessToken := new(models.OauthAccessToken)
	found = !models.OauthAccessTokenPreload(s.db).Where("token = ?", userSession.AccessToken).First(accessToken).RecordNotFound()
	if found {
		s.db.Unscoped().Where("client_id = ? AND user_id = ?", accessToken.ClientID, accessToken.UserID).Delete(models.OauthAccessToken{})
	}
}
