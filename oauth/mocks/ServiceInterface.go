package mocks

import "github.com/RichardKnop/go-oauth2-server/oauth"
import "github.com/stretchr/testify/mock"

import "net/http"
import "github.com/jinzhu/gorm"

type ServiceInterface struct {
	mock.Mock
}

func (_m *ServiceInterface) ClientExists(clientID string) bool {
	ret := _m.Called(clientID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(clientID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
func (_m *ServiceInterface) FindClientByClientID(clientID string) (*oauth.Client, error) {
	ret := _m.Called(clientID)

	var r0 *oauth.Client
	if rf, ok := ret.Get(0).(func(string) *oauth.Client); ok {
		r0 = rf(clientID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth.Client)
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
func (_m *ServiceInterface) CreateClient(clientID string, secret string, redirectURI string) (*oauth.Client, error) {
	ret := _m.Called(clientID, secret, redirectURI)

	var r0 *oauth.Client
	if rf, ok := ret.Get(0).(func(string, string, string) *oauth.Client); ok {
		r0 = rf(clientID, secret, redirectURI)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth.Client)
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
func (_m *ServiceInterface) CreateClientTx(tx *gorm.DB, clientID string, secret string, redirectURI string) (*oauth.Client, error) {
	ret := _m.Called(tx, clientID, secret, redirectURI)

	var r0 *oauth.Client
	if rf, ok := ret.Get(0).(func(*gorm.DB, string, string, string) *oauth.Client); ok {
		r0 = rf(tx, clientID, secret, redirectURI)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth.Client)
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
func (_m *ServiceInterface) AuthClient(clientID string, secret string) (*oauth.Client, error) {
	ret := _m.Called(clientID, secret)

	var r0 *oauth.Client
	if rf, ok := ret.Get(0).(func(string, string) *oauth.Client); ok {
		r0 = rf(clientID, secret)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth.Client)
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
func (_m *ServiceInterface) UserExists(username string) bool {
	ret := _m.Called(username)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(username)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
func (_m *ServiceInterface) FindUserByUsername(username string) (*oauth.User, error) {
	ret := _m.Called(username)

	var r0 *oauth.User
	if rf, ok := ret.Get(0).(func(string) *oauth.User); ok {
		r0 = rf(username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth.User)
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
func (_m *ServiceInterface) CreateUser(username string, password string) (*oauth.User, error) {
	ret := _m.Called(username, password)

	var r0 *oauth.User
	if rf, ok := ret.Get(0).(func(string, string) *oauth.User); ok {
		r0 = rf(username, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth.User)
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
func (_m *ServiceInterface) CreateUserTx(tx *gorm.DB, username string, password string) (*oauth.User, error) {
	ret := _m.Called(tx, username, password)

	var r0 *oauth.User
	if rf, ok := ret.Get(0).(func(*gorm.DB, string, string) *oauth.User); ok {
		r0 = rf(tx, username, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth.User)
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
func (_m *ServiceInterface) SetPassword(user *oauth.User, password string) error {
	ret := _m.Called(user, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(*oauth.User, string) error); ok {
		r0 = rf(user, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ServiceInterface) SetPasswordTx(tx *gorm.DB, user *oauth.User, password string) error {
	ret := _m.Called(tx, user, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, *oauth.User, string) error); ok {
		r0 = rf(tx, user, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ServiceInterface) UpdateUsername(user *oauth.User, username string) error {
	ret := _m.Called(user, username)

	var r0 error
	if rf, ok := ret.Get(0).(func(*oauth.User, string) error); ok {
		r0 = rf(user, username)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ServiceInterface) UpdateUsernameTx(db *gorm.DB, user *oauth.User, username string) error {
	ret := _m.Called(db, user, username)

	var r0 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, *oauth.User, string) error); ok {
		r0 = rf(db, user, username)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ServiceInterface) AuthUser(username string, thePassword string) (*oauth.User, error) {
	ret := _m.Called(username, thePassword)

	var r0 *oauth.User
	if rf, ok := ret.Get(0).(func(string, string) *oauth.User); ok {
		r0 = rf(username, thePassword)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth.User)
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
func (_m *ServiceInterface) GetScope(requestedScope string) (string, error) {
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
func (_m *ServiceInterface) Login(client *oauth.Client, user *oauth.User, scope string) (*oauth.AccessToken, *oauth.RefreshToken, error) {
	ret := _m.Called(client, user, scope)

	var r0 *oauth.AccessToken
	if rf, ok := ret.Get(0).(func(*oauth.Client, *oauth.User, string) *oauth.AccessToken); ok {
		r0 = rf(client, user, scope)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth.AccessToken)
		}
	}

	var r1 *oauth.RefreshToken
	if rf, ok := ret.Get(1).(func(*oauth.Client, *oauth.User, string) *oauth.RefreshToken); ok {
		r1 = rf(client, user, scope)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*oauth.RefreshToken)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*oauth.Client, *oauth.User, string) error); ok {
		r2 = rf(client, user, scope)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
func (_m *ServiceInterface) GrantAuthorizationCode(client *oauth.Client, user *oauth.User, expiresIn int, redirectURI string, scope string) (*oauth.AuthorizationCode, error) {
	ret := _m.Called(client, user, expiresIn, redirectURI, scope)

	var r0 *oauth.AuthorizationCode
	if rf, ok := ret.Get(0).(func(*oauth.Client, *oauth.User, int, string, string) *oauth.AuthorizationCode); ok {
		r0 = rf(client, user, expiresIn, redirectURI, scope)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth.AuthorizationCode)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*oauth.Client, *oauth.User, int, string, string) error); ok {
		r1 = rf(client, user, expiresIn, redirectURI, scope)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) GetValidAuthorizationCode(code string, client *oauth.Client) (*oauth.AuthorizationCode, error) {
	ret := _m.Called(code, client)

	var r0 *oauth.AuthorizationCode
	if rf, ok := ret.Get(0).(func(string, *oauth.Client) *oauth.AuthorizationCode); ok {
		r0 = rf(code, client)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth.AuthorizationCode)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *oauth.Client) error); ok {
		r1 = rf(code, client)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) GrantAccessToken(client *oauth.Client, user *oauth.User, expiresIn int, scope string) (*oauth.AccessToken, error) {
	ret := _m.Called(client, user, expiresIn, scope)

	var r0 *oauth.AccessToken
	if rf, ok := ret.Get(0).(func(*oauth.Client, *oauth.User, int, string) *oauth.AccessToken); ok {
		r0 = rf(client, user, expiresIn, scope)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth.AccessToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*oauth.Client, *oauth.User, int, string) error); ok {
		r1 = rf(client, user, expiresIn, scope)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) DeleteExpiredAccessTokens(client *oauth.Client, user *oauth.User) error {
	ret := _m.Called(client, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*oauth.Client, *oauth.User) error); ok {
		r0 = rf(client, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ServiceInterface) GetOrCreateRefreshToken(client *oauth.Client, user *oauth.User, expiresIn int, scope string) (*oauth.RefreshToken, error) {
	ret := _m.Called(client, user, expiresIn, scope)

	var r0 *oauth.RefreshToken
	if rf, ok := ret.Get(0).(func(*oauth.Client, *oauth.User, int, string) *oauth.RefreshToken); ok {
		r0 = rf(client, user, expiresIn, scope)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth.RefreshToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*oauth.Client, *oauth.User, int, string) error); ok {
		r1 = rf(client, user, expiresIn, scope)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) GetValidRefreshToken(token string, client *oauth.Client) (*oauth.RefreshToken, error) {
	ret := _m.Called(token, client)

	var r0 *oauth.RefreshToken
	if rf, ok := ret.Get(0).(func(string, *oauth.Client) *oauth.RefreshToken); ok {
		r0 = rf(token, client)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth.RefreshToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *oauth.Client) error); ok {
		r1 = rf(token, client)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) Authenticate(token string) (*oauth.AccessToken, error) {
	ret := _m.Called(token)

	var r0 *oauth.AccessToken
	if rf, ok := ret.Get(0).(func(string) *oauth.AccessToken); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth.AccessToken)
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
func (_m *ServiceInterface) TokensHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}
func (_m *ServiceInterface) IntrospectHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}
