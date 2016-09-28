package oauth

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/RichardKnop/go-oauth2-server/routes"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	GetConfig() *config.Config
	RestrictToRoles(allowedRoles ...string)
	IsRoleAllowed(role string) bool
	GetRoutes() []routes.Route
	RegisterRoutes(router *mux.Router, prefix string)
	ClientExists(clientID string) bool
	FindClientByClientID(clientID string) (*Client, error)
	CreateClient(clientID, secret, redirectURI string) (*Client, error)
	CreateClientTx(tx *gorm.DB, clientID, secret, redirectURI string) (*Client, error)
	AuthClient(clientID, secret string) (*Client, error)
	UserExists(username string) bool
	FindUserByUsername(username string) (*User, error)
	CreateUser(roleID, username, password string) (*User, error)
	CreateUserTx(tx *gorm.DB, roleID, username, password string) (*User, error)
	SetPassword(user *User, password string) error
	SetPasswordTx(tx *gorm.DB, user *User, password string) error
	UpdateUsername(user *User, username string) error
	UpdateUsernameTx(db *gorm.DB, user *User, username string) error
	AuthUser(username, thePassword string) (*User, error)
	GetScope(requestedScope string) (string, error)
	Login(client *Client, user *User, scope string) (*AccessToken, *RefreshToken, error)
	GrantAuthorizationCode(client *Client, user *User, expiresIn int, redirectURI, scope string) (*AuthorizationCode, error)
	GrantAccessToken(client *Client, user *User, expiresIn int, scope string) (*AccessToken, error)
	GetOrCreateRefreshToken(client *Client, user *User, expiresIn int, scope string) (*RefreshToken, error)
	GetValidRefreshToken(token string, client *Client) (*RefreshToken, error)
	Authenticate(token string) (*AccessToken, error)
	NewIntrospectResponseFromAccessToken(accessToken *AccessToken) (*IntrospectResponse, error)
	NewIntrospectResponseFromRefreshToken(refreshToken *RefreshToken) (*IntrospectResponse, error)
}
