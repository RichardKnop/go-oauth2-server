package oauth_test

import (
	"time"

	"github.com/RichardKnop/go-oauth2-server/models"
	"github.com/RichardKnop/uuid"
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestGrantAccessToken() {
	var (
		accessToken *models.OauthAccessToken
		err         error
		tokens      []*models.OauthAccessToken
	)

	// Grant a client only access token
	accessToken, err = suite.service.GrantAccessToken(
		suite.clients[0],       // client
		nil,                    // user
		3600,                   // expires in
		"scope doesn't matter", // scope
	)

	// Error should be Nil
	assert.Nil(suite.T(), err)

	// Correct access token object should be returned
	if assert.NotNil(suite.T(), accessToken) {
		// Fetch all access tokens
		models.OauthAccessTokenPreload(suite.db).Order("created_at").Find(&tokens)

		// There should be just one right now
		assert.Equal(suite.T(), 1, len(tokens))

		// And the token should match the one returned by the grant method
		assert.Equal(suite.T(), tokens[0].Token, accessToken.Token)

		// Client id should be set
		assert.True(suite.T(), tokens[0].ClientID.Valid)
		assert.Equal(suite.T(), string(suite.clients[0].ID), tokens[0].ClientID.String)

		// User id should be nil
		assert.False(suite.T(), tokens[0].UserID.Valid)
	}

	// Grant a user specific access token
	accessToken, err = suite.service.GrantAccessToken(
		suite.clients[0],       // client
		suite.users[0],         // user
		3600,                   // expires in
		"scope doesn't matter", // scope
	)

	// Error should be Nil
	assert.Nil(suite.T(), err)

	// Correct access token object should be returned
	if assert.NotNil(suite.T(), accessToken) {
		// Fetch all access tokens
		models.OauthAccessTokenPreload(suite.db).Order("created_at").Find(&tokens)

		// There should be 2 tokens now
		assert.Equal(suite.T(), 2, len(tokens))

		// And the second token should match the one returned by the grant method
		assert.Equal(suite.T(), tokens[1].Token, accessToken.Token)

		// Client id should be set
		assert.True(suite.T(), tokens[1].ClientID.Valid)
		assert.Equal(suite.T(), string(suite.clients[0].ID), tokens[1].ClientID.String)

		// User id should be set
		assert.True(suite.T(), tokens[1].UserID.Valid)
		assert.Equal(suite.T(), string(suite.users[0].ID), tokens[1].UserID.String)
	}
}

func (suite *OauthTestSuite) TestGrantAccessTokenDeletesExpiredTokens() {
	var (
		testAccessTokens = []*models.OauthAccessToken{
			// Expired access token with a user
			{
				MyGormModel: models.MyGormModel{
					ID:        uuid.New(),
					CreatedAt: time.Now().UTC(),
				},
				Token:     "test_token_1",
				ExpiresAt: time.Now().UTC().Add(-10 * time.Second),
				Client:    suite.clients[0],
				User:      suite.users[0],
			},
			// Expired access token without a user
			{
				MyGormModel: models.MyGormModel{
					ID:        uuid.New(),
					CreatedAt: time.Now().UTC(),
				},
				Token:     "test_token_2",
				ExpiresAt: time.Now().UTC().Add(-10 * time.Second),
				Client:    suite.clients[0],
			},
			// Access token with a user
			{
				MyGormModel: models.MyGormModel{
					ID:        uuid.New(),
					CreatedAt: time.Now().UTC(),
				},
				Token:     "test_token_3",
				ExpiresAt: time.Now().UTC().Add(+10 * time.Second),
				Client:    suite.clients[0],
				User:      suite.users[0],
			},
			// Access token without a user
			{
				MyGormModel: models.MyGormModel{
					ID:        uuid.New(),
					CreatedAt: time.Now().UTC(),
				},
				Token:     "test_token_4",
				ExpiresAt: time.Now().UTC().Add(+10 * time.Second),
				Client:    suite.clients[0],
			},
		}
		err            error
		notFound       bool
		existingTokens []string
	)

	// Insert test access tokens
	for _, testAccessToken := range testAccessTokens {
		err = suite.db.Create(testAccessToken).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	// This should only delete test_token_1
	_, err = suite.service.GrantAccessToken(
		suite.clients[0],       // client
		suite.users[0],         // user
		3600,                   // expires in
		"scope doesn't matter", // scope
	)
	assert.NoError(suite.T(), err)

	// Check the test_token_1 was deleted
	notFound = suite.db.Unscoped().Where("token = ?", "test_token_1").
		First(new(models.OauthAccessToken)).RecordNotFound()
	assert.True(suite.T(), notFound)

	// Check the other three tokens are still around
	existingTokens = []string{
		"test_token_2",
		"test_token_3",
		"test_token_4",
	}
	for _, token := range existingTokens {
		notFound = suite.db.Unscoped().Where("token = ?", token).
			First(new(models.OauthAccessToken)).RecordNotFound()
		assert.False(suite.T(), notFound)
	}

	// This should only delete test_token_2
	_, err = suite.service.GrantAccessToken(
		suite.clients[0],       // client
		nil,                    // user
		3600,                   // expires in
		"scope doesn't matter", // scope
	)
	assert.NoError(suite.T(), err)

	// Check the test_token_2 was deleted
	notFound = suite.db.Unscoped().Where("token = ?", "test_token_2").
		First(new(models.OauthAccessToken)).RecordNotFound()
	assert.True(suite.T(), notFound)

	// Check that last two tokens are still around
	existingTokens = []string{
		"test_token_3",
		"test_token_4",
	}
	for _, token := range existingTokens {
		notFound := suite.db.Unscoped().Where("token = ?", token).
			First(new(models.OauthAccessToken)).RecordNotFound()
		assert.False(suite.T(), notFound)
	}
}
