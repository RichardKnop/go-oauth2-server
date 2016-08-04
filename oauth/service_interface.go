package oauth

import (
	"net/http"

	"github.com/jinzhu/gorm"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	ClientExists(clientID string) bool
	FindClientByClientID(clientID string) (*Client, error)
	CreateClient(clientID, secret, redirectURI string) (*Client, error)
	CreateClientTx(tx *gorm.DB, clientID, secret, redirectURI string) (*Client, error)
	AuthClient(clientID, secret string) (*Client, error)
	UserExists(username string) bool
	FindUserByUsername(username string) (*User, error)
	CreateUser(username, password string) (*User, error)
	CreateUserTx(tx *gorm.DB, username, password string) (*User, error)
	SetPassword(user *User, password string) error
	SetPasswordTx(tx *gorm.DB, user *User, password string) error
	UpdateUsername(user *User, username string) error
	UpdateUsernameTx(db *gorm.DB, user *User, username string) error
	AuthUser(username, thePassword string) (*User, error)
	GetScope(requestedScope string) (string, error)
	Login(client *Client, user *User, scope string) (*AccessToken, *RefreshToken, error)
	GrantAuthorizationCode(client *Client, user *User, expiresIn int, redirectURI, scope string) (*AuthorizationCode, error)
	GetValidAuthorizationCode(code string, client *Client) (*AuthorizationCode, error)
	GrantAccessToken(client *Client, user *User, expiresIn int, scope string) (*AccessToken, error)
	DeleteExpiredAccessTokens(client *Client, user *User) error
	GetOrCreateRefreshToken(client *Client, user *User, expiresIn int, scope string) (*RefreshToken, error)
	GetValidRefreshToken(token string, client *Client) (*RefreshToken, error)
	Authenticate(token string) (*AccessToken, error)

	// Needed for the NewRoutes to be able to register handlers
	TokensHandler(w http.ResponseWriter, r *http.Request)
	IntrospectHandler(w http.ResponseWriter, r *http.Request)
}
