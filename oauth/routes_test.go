package oauth_test

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestTokensRouteIsValid() {
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/oauth/tokens",
		nil,
	)
	assert.NoError(suite.T(), err, "New request should not cause an error")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route, "Expected to find a route match") {
		assert.Equal(suite.T(), "oauth_tokens", match.Route.GetName(), "Expected route to be matched")
	}
}

func (suite *OauthTestSuite) TestIntrospectRouteIsValid() {
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/oauth/introspect",
		nil,
	)
	assert.NoError(suite.T(), err, "New request should not cause an error")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route, "Expected to find a route match") {
		assert.Equal(suite.T(), "oauth_introspect", match.Route.GetName(), "Expected route to be matched")
	}
}
