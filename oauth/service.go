package oauth

import (
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/jinzhu/gorm"
)

// type ServiceInterface interface {
// 	AuthClient(clientID, secret string) (*Client, error)
// 	CreateClient(clientID, secret string) (*Client, error)
// 	FindClientByClientID(clientID string) (*Client, error)
// 	UserExists(username string) bool
// 	AuthUser(username, password string) (*User, error)
// 	CreateUser(username, password string) (*User, error)
// 	GetScope(requestedScope string) (string, error)
//  GrantAuthorizationCode(client *Client, user *User, scope string) (*AuthorizationCode, error)
// 	GrantAccessToken(client *Client, user *User, scope string) (*AccessToken, error)
// 	GetOrCreateRefreshToken(client *Client, user *User, scope string) (*RefreshToken, error)
// 	ValidateRefreshToken(token string, client *Client) (*RefreshToken, error)
// 	Authenticate(token string) error
// }

// Service struct keeps config and db objects to avoid passing them around
type Service struct {
	cnf *config.Config
	db  *gorm.DB
}

var theService *Service

// NewService starts a new Service instance
func NewService(cnf *config.Config, db *gorm.DB) *Service {
	theService = &Service{cnf: cnf, db: db}
	return theService
}

// GetService returns internal Service instance
func GetService() *Service {
	return theService
}
