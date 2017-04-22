package mocks

import "github.com/RichardKnop/go-oauth2-server/oauth"
import "github.com/stretchr/testify/mock"

import "github.com/RichardKnop/go-oauth2-server/config"
import "github.com/RichardKnop/go-oauth2-server/models"
import "github.com/RichardKnop/go-oauth2-server/util/routes"
import "github.com/gorilla/mux"
import "github.com/jinzhu/gorm"

type ServiceInterface struct {
	mock.Mock
}

func (_m *ServiceInterface) GetConfig() *config.Config {
	ret := _m.Called()

	var r0 *config.Config
	if rf, ok := ret.Get(0).(func() *config.Config); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*config.Config)
		}
	}

	return r0
}
func (_m *ServiceInterface) RestrictToRoles(allowedRoles ...string) {
	_m.Called(allowedRoles)
}
func (_m *ServiceInterface) IsRoleAllowed(role string) bool {
	ret := _m.Called(role)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(role)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
func (_m *ServiceInterface) GetRoutes() []routes.Route {
	ret := _m.Called()

	var r0 []routes.Route
	if rf, ok := ret.Get(0).(func() []routes.Route); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]routes.Route)
		}
	}

	return r0
}
func (_m *ServiceInterface) RegisterRoutes(router *mux.Router, prefix string) {
	_m.Called(router, prefix)
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
func (_m *ServiceInterface) FindClientByClientID(clientID string) (*models.OauthClient, error) {
	ret := _m.Called(clientID)

	var r0 *models.OauthClient
	if rf, ok := ret.Get(0).(func(string) *models.OauthClient); ok {
		r0 = rf(clientID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.OauthClient)
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
func (_m *ServiceInterface) CreateClient(clientID string, secret string, redirectURI string) (*models.OauthClient, error) {
	ret := _m.Called(clientID, secret, redirectURI)

	var r0 *models.OauthClient
	if rf, ok := ret.Get(0).(func(string, string, string) *models.OauthClient); ok {
		r0 = rf(clientID, secret, redirectURI)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.OauthClient)
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
func (_m *ServiceInterface) CreateClientTx(tx *gorm.DB, clientID string, secret string, redirectURI string) (*models.OauthClient, error) {
	ret := _m.Called(tx, clientID, secret, redirectURI)

	var r0 *models.OauthClient
	if rf, ok := ret.Get(0).(func(*gorm.DB, string, string, string) *models.OauthClient); ok {
		r0 = rf(tx, clientID, secret, redirectURI)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.OauthClient)
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
func (_m *ServiceInterface) AuthClient(clientID string, secret string) (*models.OauthClient, error) {
	ret := _m.Called(clientID, secret)

	var r0 *models.OauthClient
	if rf, ok := ret.Get(0).(func(string, string) *models.OauthClient); ok {
		r0 = rf(clientID, secret)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.OauthClient)
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
func (_m *ServiceInterface) FindUserByUsername(username string) (*models.OauthUser, error) {
	ret := _m.Called(username)

	var r0 *models.OauthUser
	if rf, ok := ret.Get(0).(func(string) *models.OauthUser); ok {
		r0 = rf(username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.OauthUser)
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
func (_m *ServiceInterface) CreateUser(roleID string, username string, password string) (*models.OauthUser, error) {
	ret := _m.Called(roleID, username, password)

	var r0 *models.OauthUser
	if rf, ok := ret.Get(0).(func(string, string, string) *models.OauthUser); ok {
		r0 = rf(roleID, username, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.OauthUser)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(roleID, username, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) CreateUserTx(tx *gorm.DB, roleID string, username string, password string) (*models.OauthUser, error) {
	ret := _m.Called(tx, roleID, username, password)

	var r0 *models.OauthUser
	if rf, ok := ret.Get(0).(func(*gorm.DB, string, string, string) *models.OauthUser); ok {
		r0 = rf(tx, roleID, username, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.OauthUser)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB, string, string, string) error); ok {
		r1 = rf(tx, roleID, username, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) SetPassword(user *models.OauthUser, password string) error {
	ret := _m.Called(user, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.OauthUser, string) error); ok {
		r0 = rf(user, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ServiceInterface) SetPasswordTx(tx *gorm.DB, user *models.OauthUser, password string) error {
	ret := _m.Called(tx, user, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, *models.OauthUser, string) error); ok {
		r0 = rf(tx, user, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ServiceInterface) UpdateUsername(user *models.OauthUser, username string) error {
	ret := _m.Called(user, username)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.OauthUser, string) error); ok {
		r0 = rf(user, username)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ServiceInterface) UpdateUsernameTx(db *gorm.DB, user *models.OauthUser, username string) error {
	ret := _m.Called(db, user, username)

	var r0 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, *models.OauthUser, string) error); ok {
		r0 = rf(db, user, username)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ServiceInterface) AuthUser(username string, thePassword string) (*models.OauthUser, error) {
	ret := _m.Called(username, thePassword)

	var r0 *models.OauthUser
	if rf, ok := ret.Get(0).(func(string, string) *models.OauthUser); ok {
		r0 = rf(username, thePassword)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.OauthUser)
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
func (_m *ServiceInterface) Login(client *models.OauthClient, user *models.OauthUser, scope string) (*models.OauthAccessToken, *models.OauthRefreshToken, error) {
	ret := _m.Called(client, user, scope)

	var r0 *models.OauthAccessToken
	if rf, ok := ret.Get(0).(func(*models.OauthClient, *models.OauthUser, string) *models.OauthAccessToken); ok {
		r0 = rf(client, user, scope)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.OauthAccessToken)
		}
	}

	var r1 *models.OauthRefreshToken
	if rf, ok := ret.Get(1).(func(*models.OauthClient, *models.OauthUser, string) *models.OauthRefreshToken); ok {
		r1 = rf(client, user, scope)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*models.OauthRefreshToken)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*models.OauthClient, *models.OauthUser, string) error); ok {
		r2 = rf(client, user, scope)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
func (_m *ServiceInterface) GrantAuthorizationCode(client *models.OauthClient, user *models.OauthUser, expiresIn int, redirectURI string, scope string) (*models.OauthAuthorizationCode, error) {
	ret := _m.Called(client, user, expiresIn, redirectURI, scope)

	var r0 *models.OauthAuthorizationCode
	if rf, ok := ret.Get(0).(func(*models.OauthClient, *models.OauthUser, int, string, string) *models.OauthAuthorizationCode); ok {
		r0 = rf(client, user, expiresIn, redirectURI, scope)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.OauthAuthorizationCode)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.OauthClient, *models.OauthUser, int, string, string) error); ok {
		r1 = rf(client, user, expiresIn, redirectURI, scope)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) GrantAccessToken(client *models.OauthClient, user *models.OauthUser, expiresIn int, scope string) (*models.OauthAccessToken, error) {
	ret := _m.Called(client, user, expiresIn, scope)

	var r0 *models.OauthAccessToken
	if rf, ok := ret.Get(0).(func(*models.OauthClient, *models.OauthUser, int, string) *models.OauthAccessToken); ok {
		r0 = rf(client, user, expiresIn, scope)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.OauthAccessToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.OauthClient, *models.OauthUser, int, string) error); ok {
		r1 = rf(client, user, expiresIn, scope)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) GetOrCreateRefreshToken(client *models.OauthClient, user *models.OauthUser, expiresIn int, scope string) (*models.OauthRefreshToken, error) {
	ret := _m.Called(client, user, expiresIn, scope)

	var r0 *models.OauthRefreshToken
	if rf, ok := ret.Get(0).(func(*models.OauthClient, *models.OauthUser, int, string) *models.OauthRefreshToken); ok {
		r0 = rf(client, user, expiresIn, scope)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.OauthRefreshToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.OauthClient, *models.OauthUser, int, string) error); ok {
		r1 = rf(client, user, expiresIn, scope)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) GetValidRefreshToken(token string, client *models.OauthClient) (*models.OauthRefreshToken, error) {
	ret := _m.Called(token, client)

	var r0 *models.OauthRefreshToken
	if rf, ok := ret.Get(0).(func(string, *models.OauthClient) *models.OauthRefreshToken); ok {
		r0 = rf(token, client)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.OauthRefreshToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *models.OauthClient) error); ok {
		r1 = rf(token, client)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) Authenticate(token string) (*models.OauthAccessToken, error) {
	ret := _m.Called(token)

	var r0 *models.OauthAccessToken
	if rf, ok := ret.Get(0).(func(string) *models.OauthAccessToken); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.OauthAccessToken)
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
func (_m *ServiceInterface) NewIntrospectResponseFromAccessToken(accessToken *models.OauthAccessToken) (*oauth.IntrospectResponse, error) {
	ret := _m.Called(accessToken)

	var r0 *oauth.IntrospectResponse
	if rf, ok := ret.Get(0).(func(*models.OauthAccessToken) *oauth.IntrospectResponse); ok {
		r0 = rf(accessToken)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth.IntrospectResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.OauthAccessToken) error); ok {
		r1 = rf(accessToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) NewIntrospectResponseFromRefreshToken(refreshToken *models.OauthRefreshToken) (*oauth.IntrospectResponse, error) {
	ret := _m.Called(refreshToken)

	var r0 *oauth.IntrospectResponse
	if rf, ok := ret.Get(0).(func(*models.OauthRefreshToken) *oauth.IntrospectResponse); ok {
		r0 = rf(refreshToken)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth.IntrospectResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.OauthRefreshToken) error); ok {
		r1 = rf(refreshToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
