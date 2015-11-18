package oauth

// ServiceInterface defines exported methods
type ServiceInterface interface {
	ClientExists(clientID string) bool
	FindClientByClientID(clientID string) (*Client, error)
	CreateClient(clientID, secret, redirectURI string) (*Client, error)
	AuthClient(clientID, secret string) (*Client, error)
	UserExists(username string) bool
	FindUserByUsername(username string) (*User, error)
	CreateUser(username, thePassword string) (*User, error)
	AuthUser(username, thePassword string) (*User, error)
	GetScope(requestedScope string) (string, error)
	GrantAuthorizationCode(client *Client, user *User, redirectURI, scope string) (*AuthorizationCode, error)
	GrantAccessToken(client *Client, user *User, scope string) (*AccessToken, error)
	GetOrCreateRefreshToken(client *Client, user *User, scope string) (*RefreshToken, error)
	GetValidRefreshToken(token string, client *Client) (*RefreshToken, error)
	Authenticate(token string) error
}
