package oauth

import (
	"time"

	"github.com/RichardKnop/go-oauth2-server/models"
)

// GrantAccessToken deletes old tokens and grants a new access token
func (s *Service) GrantAccessToken(client *models.OauthClient, user *models.OauthUser, expiresIn int, scope string) (*models.OauthAccessToken, error) {
	// Begin a transaction
	tx := s.db.Begin()

	// Delete expired access tokens
	query := tx.Unscoped().Where("client_id = ?", client.ID)
	if user != nil && len([]rune(user.ID)) > 0 {
		query = query.Where("user_id = ?", user.ID)
	} else {
		query = query.Where("user_id IS NULL")
	}
	if err := query.Where("expires_at <= ?", time.Now()).Delete(new(models.OauthAccessToken)).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Create a new access token
	accessToken := models.NewOauthAccessToken(client, user, expiresIn, scope)
	if err := tx.Create(accessToken).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}
	accessToken.Client = client
	accessToken.User = user

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	return accessToken, nil
}
