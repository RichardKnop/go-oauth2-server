package oauth_test

import (
	"github.com/RichardKnop/go-oauth2-server/models"
	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/RichardKnop/go-oauth2-server/oauth/roles"
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestFindRoleByID() {
	var (
		role *models.OauthRole
		err  error
	)

	// Let's try to find a role by a bogus ID
	role, err = suite.service.FindRoleByID("bogus")

	// Role should be nil
	assert.Nil(suite.T(), role)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), oauth.ErrRoleNotFound, err)
	}

	// Now let's pass a valid ID
	role, err = suite.service.FindRoleByID(roles.User)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct role should be returned
	if assert.NotNil(suite.T(), role) {
		assert.Equal(suite.T(), roles.User, role.ID)
	}
}
