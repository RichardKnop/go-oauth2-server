package oauth

import (
	"github.com/adam-hanna/go-oauth2-server/config"
	"github.com/adam-hanna/go-oauth2-server/models"
	"github.com/adam-hanna/go-oauth2-server/oauth/roles"
	"github.com/adam-hanna/go-oauth2-server/session"
	"github.com/jinzhu/gorm"
)

// Service struct keeps objects to avoid passing them around
type Service struct {
	cnf          *config.Config
	db           *gorm.DB
	allowedRoles []string
}

// NewService returns a new Service instance
func NewService(cnf *config.Config, db *gorm.DB) *Service {
	return &Service{
		cnf:          cnf,
		db:           db,
		allowedRoles: []string{roles.Superuser, roles.User},
	}
}

// GetConfig returns config.Config instance
func (s *Service) GetConfig() *config.Config {
	return s.cnf
}

// RestrictToRoles restricts this service to only specified roles
func (s *Service) RestrictToRoles(allowedRoles ...string) {
	s.allowedRoles = allowedRoles
}

// IsRoleAllowed returns true if the role is allowed to use this service
func (s *Service) IsRoleAllowed(role string) bool {
	for _, allowedRole := range s.allowedRoles {
		if role == allowedRole {
			return true
		}
	}
	return false
}

// Close stops any running services
func (s *Service) Close() {}

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
