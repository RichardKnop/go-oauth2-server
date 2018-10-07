package oauth_test

import (
	"github.com/RichardKnop/go-oauth2-server/models"
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestGrantAuthorizationCode() {
	var (
		authorizationCode *models.OauthAuthorizationCode
		err               error
		codes             []*models.OauthAuthorizationCode
	)

	// Grant an authorization code
	authorizationCode, err = suite.service.GrantAuthorizationCode(
		suite.clients[0],              // client
		suite.users[0],                // user
		3600,                          // expires in
		"redirect URI doesn't matter", // redirect URI
		"scope doesn't matter",        // scope
	)

	// Error should be Nil
	assert.Nil(suite.T(), err)

	// Correct authorization code object should be returned
	if assert.NotNil(suite.T(), authorizationCode) {
		// Fetch all auth codes
		models.OauthAuthorizationCodePreload(suite.db).Order("created_at").Find(&codes)

		// There should be just one right now
		assert.Equal(suite.T(), 1, len(codes))

		// And the code should match the one returned by the grant method
		assert.Equal(suite.T(), codes[0].Code, authorizationCode.Code)

		// Client ID should be set
		assert.True(suite.T(), codes[0].ClientID.Valid)
		assert.Equal(suite.T(), string(suite.clients[0].ID), codes[0].ClientID.String)

		// User ID should be set
		assert.True(suite.T(), codes[0].UserID.Valid)
		assert.Equal(suite.T(), string(suite.users[0].ID), codes[0].UserID.String)
	}
}
