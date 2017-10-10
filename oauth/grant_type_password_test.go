package oauth_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/RichardKnop/go-oauth2-server/models"
	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/RichardKnop/go-oauth2-server/oauth/roles"
	"github.com/RichardKnop/go-oauth2-server/oauth/tokentypes"
	"github.com/RichardKnop/go-oauth2-server/test-util"
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestPasswordGrant() {
	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/tokens", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")
	r.PostForm = url.Values{
		"grant_type": {"password"},
		"username":   {"test@user"},
		"password":   {"test_password"},
		"scope":      {"read_write"},
	}

	// Serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Fetch data
	accessToken, refreshToken := new(models.OauthAccessToken), new(models.OauthRefreshToken)
	assert.False(suite.T(), models.OauthAccessTokenPreload(suite.db).
		Last(accessToken).RecordNotFound())
	assert.False(suite.T(), models.OauthRefreshTokenPreload(suite.db).
		Last(refreshToken).RecordNotFound())

	// Check the response
	expected := &oauth.AccessTokenResponse{
		UserID:       accessToken.UserID.String,
		AccessToken:  accessToken.Token,
		ExpiresIn:    3600,
		TokenType:    tokentypes.Bearer,
		Scope:        "read_write",
		RefreshToken: refreshToken.Token,
	}
	testutil.TestResponseObject(suite.T(), w, expected, 200)
}

func (suite *OauthTestSuite) TestPasswordGrantWithRoleRestriction() {
	suite.service.RestrictToRoles(roles.Superuser)

	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/tokens", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")
	r.PostForm = url.Values{
		"grant_type": {"password"},
		"username":   {"test@user"},
		"password":   {"test_password"},
		"scope":      {"read_write"},
	}

	// Serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the response
	testutil.TestResponseForError(
		suite.T(),
		w,
		oauth.ErrInvalidUsernameOrPassword.Error(),
		401,
	)

	suite.service.RestrictToRoles(roles.Superuser, roles.User)
}
