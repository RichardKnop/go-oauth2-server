package oauth

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

// ServiceMock is a mocked object implementing ServiceInterface
type ServiceMock struct {
	mock.Mock
}

// ClientExists ...
func (_m *ServiceMock) ClientExists(clientID string) bool {
	ret := _m.Called(clientID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(clientID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// FindClientByClientID ...
func (_m *ServiceMock) FindClientByClientID(clientID string) (*Client, error) {
	ret := _m.Called(clientID)

	var r0 *Client
	if rf, ok := ret.Get(0).(func(string) *Client); ok {
		r0 = rf(clientID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Client)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(clientID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateClient ...
func (_m *ServiceMock) CreateClient(clientID string, secret string, redirectURI string) (*Client, error) {
	ret := _m.Called(clientID, secret, redirectURI)

	var r0 *Client
	if rf, ok := ret.Get(0).(func(string, string, string) *Client); ok {
		r0 = rf(clientID, secret, redirectURI)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Client)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(clientID, secret, redirectURI)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateClientTx ...
func (_m *ServiceMock) CreateClientTx(tx *gorm.DB, clientID string, secret string, redirectURI string) (*Client, error) {
	ret := _m.Called(tx, clientID, secret, redirectURI)

	var r0 *Client
	if rf, ok := ret.Get(0).(func(*gorm.DB, string, string, string) *Client); ok {
		r0 = rf(tx, clientID, secret, redirectURI)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Client)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB, string, string, string) error); ok {
		r1 = rf(tx, clientID, secret, redirectURI)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthClient ...
func (_m *ServiceMock) AuthClient(clientID string, secret string) (*Client, error) {
	ret := _m.Called(clientID, secret)

	var r0 *Client
	if rf, ok := ret.Get(0).(func(string, string) *Client); ok {
		r0 = rf(clientID, secret)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Client)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(clientID, secret)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserExists ...
func (_m *ServiceMock) UserExists(username string) bool {
	ret := _m.Called(username)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(username)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// FindUserByUsername ...
func (_m *ServiceMock) FindUserByUsername(username string) (*User, error) {
	ret := _m.Called(username)

	var r0 *User
	if rf, ok := ret.Get(0).(func(string) *User); ok {
		r0 = rf(username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUser ...
func (_m *ServiceMock) CreateUser(username string, password string) (*User, error) {
	ret := _m.Called(username, password)

	var r0 *User
	if rf, ok := ret.Get(0).(func(string, string) *User); ok {
		r0 = rf(username, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(username, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUserTx ...
func (_m *ServiceMock) CreateUserTx(tx *gorm.DB, username string, password string) (*User, error) {
	ret := _m.Called(tx, username, password)

	var r0 *User
	if rf, ok := ret.Get(0).(func(*gorm.DB, string, string) *User); ok {
		r0 = rf(tx, username, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB, string, string) error); ok {
		r1 = rf(tx, username, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetPassword ...
func (_m *ServiceMock) SetPassword(user *User, password string) error {
	ret := _m.Called(user, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(*User, string) error); ok {
		r0 = rf(user, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AuthUser ...
func (_m *ServiceMock) AuthUser(username string, thePassword string) (*User, error) {
	ret := _m.Called(username, thePassword)

	var r0 *User
	if rf, ok := ret.Get(0).(func(string, string) *User); ok {
		r0 = rf(username, thePassword)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(username, thePassword)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetScope ...
func (_m *ServiceMock) GetScope(requestedScope string) (string, error) {
	ret := _m.Called(requestedScope)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(requestedScope)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(requestedScope)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GrantAuthorizationCode ...
func (_m *ServiceMock) GrantAuthorizationCode(client *Client, user *User, expiresIn int, redirectURI string, scope string) (*AuthorizationCode, error) {
	ret := _m.Called(client, user, expiresIn, redirectURI, scope)

	var r0 *AuthorizationCode
	if rf, ok := ret.Get(0).(func(*Client, *User, int, string, string) *AuthorizationCode); ok {
		r0 = rf(client, user, expiresIn, redirectURI, scope)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*AuthorizationCode)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*Client, *User, int, string, string) error); ok {
		r1 = rf(client, user, expiresIn, redirectURI, scope)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GrantAccessToken ...
func (_m *ServiceMock) GrantAccessToken(client *Client, user *User, expiresIn int, scope string) (*AccessToken, error) {
	ret := _m.Called(client, user, expiresIn, scope)

	var r0 *AccessToken
	if rf, ok := ret.Get(0).(func(*Client, *User, int, string) *AccessToken); ok {
		r0 = rf(client, user, expiresIn, scope)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*AccessToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*Client, *User, int, string) error); ok {
		r1 = rf(client, user, expiresIn, scope)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrCreateRefreshToken ...
func (_m *ServiceMock) GetOrCreateRefreshToken(client *Client, user *User, expiresIn int, scope string) (*RefreshToken, error) {
	ret := _m.Called(client, user, expiresIn, scope)

	var r0 *RefreshToken
	if rf, ok := ret.Get(0).(func(*Client, *User, int, string) *RefreshToken); ok {
		r0 = rf(client, user, expiresIn, scope)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*RefreshToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*Client, *User, int, string) error); ok {
		r1 = rf(client, user, expiresIn, scope)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetValidRefreshToken ...
func (_m *ServiceMock) GetValidRefreshToken(token string, client *Client) (*RefreshToken, error) {
	ret := _m.Called(token, client)

	var r0 *RefreshToken
	if rf, ok := ret.Get(0).(func(string, *Client) *RefreshToken); ok {
		r0 = rf(token, client)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*RefreshToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *Client) error); ok {
		r1 = rf(token, client)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Authenticate ...
func (_m *ServiceMock) Authenticate(token string) (*AccessToken, error) {
	ret := _m.Called(token)

	var r0 *AccessToken
	if rf, ok := ret.Get(0).(func(string) *AccessToken); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*AccessToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *ServiceMock) tokensHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

func (_m *ServiceMock) introspectHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}
