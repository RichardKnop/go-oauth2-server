package oauth

import (
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestAuthenticate() {
	var (
		accessToken *AccessToken
		err         error
	)

	// Insert some test access tokens
	testAccessTokens := []*AccessToken{
		// Expired access token
		&AccessToken{
			Token:     "test_expired_token",
			ExpiresAt: time.Now().Add(-10 * time.Second),
			Client:    suite.clients[0],
			User:      suite.users[0],
		},
		// Access token without a user
		&AccessToken{
			Token:     "test_client_token",
			ExpiresAt: time.Now().Add(+10 * time.Second),
			Client:    suite.clients[0],
		},
		// Access token with a user
		&AccessToken{
			Token:     "test_user_token",
			ExpiresAt: time.Now().Add(+10 * time.Second),
			Client:    suite.clients[0],
			User:      suite.users[0],
		},
	}
	for _, testAccessToken := range testAccessTokens {
		err := suite.db.Create(testAccessToken).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	// Test passing an empty token
	accessToken, err = suite.service.Authenticate("")

	// Access token should be nil
	assert.Nil(suite.T(), accessToken)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrAccessTokenNotFound, err)
	}

	// Test passing a bogus token
	accessToken, err = suite.service.Authenticate("bogus")

	// Access token should be nil
	assert.Nil(suite.T(), accessToken)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrAccessTokenNotFound, err)
	}

	// Test passing an expired token
	accessToken, err = suite.service.Authenticate("test_expired_token")

	// Access token should be nil
	assert.Nil(suite.T(), accessToken)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrAccessTokenExpired, err)
	}

	// Test passing a valid client token
	accessToken, err = suite.service.Authenticate("test_client_token")

	// Correct access token should be returned
	if assert.NotNil(suite.T(), accessToken) {
		assert.Equal(suite.T(), "test_client_token", accessToken.Token)
		assert.EqualValues(suite.T(), suite.clients[0].ID, accessToken.ClientID.Int64)
		assert.False(suite.T(), accessToken.UserID.Valid)
	}

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Test passing a valid user token
	accessToken, err = suite.service.Authenticate("test_user_token")

	// Correct access token should be returned
	if assert.NotNil(suite.T(), accessToken) {
		assert.Equal(suite.T(), "test_user_token", accessToken.Token)
		assert.EqualValues(suite.T(), suite.clients[0].ID, accessToken.ClientID.Int64)
		assert.EqualValues(suite.T(), suite.users[0].ID, accessToken.UserID.Int64)
	}

	// Error should be nil
	assert.Nil(suite.T(), err)
}

func (suite *OauthTestSuite) TestAuthenticateRollingRefreshToken() {
	var (
		testAccessTokens  []*AccessToken
		testRefreshTokens []*RefreshToken
		accessToken       *AccessToken
		err               error
		refreshTokens     []*RefreshToken
	)

	// Insert some test access tokens
	testAccessTokens = []*AccessToken{
		&AccessToken{
			Token:     "test_token_1",
			ExpiresAt: time.Now().Add(+10 * time.Second),
			Client:    suite.clients[0],
			User:      suite.users[0],
		},
		&AccessToken{
			Token:     "test_token_2",
			ExpiresAt: time.Now().Add(+10 * time.Second),
			Client:    suite.clients[0],
		},
		&AccessToken{
			Token:     "test_token_3",
			ExpiresAt: time.Now().Add(+10 * time.Second),
			Client:    suite.clients[0],
			User:      suite.users[1],
		},
	}
	for _, testAccessToken := range testAccessTokens {
		err = suite.db.Create(testAccessToken).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	// Insert some test access tokens
	testRefreshTokens = []*RefreshToken{
		&RefreshToken{
			Token:     "test_token_1",
			ExpiresAt: time.Now().Add(+10 * time.Second),
			Client:    suite.clients[0],
			User:      suite.users[0],
		},
		&RefreshToken{
			Token:     "test_token_2",
			ExpiresAt: time.Now().Add(+10 * time.Second),
			Client:    suite.clients[0],
		},
		&RefreshToken{
			Token:     "test_token_3",
			ExpiresAt: time.Now().Add(+10 * time.Second),
			Client:    suite.clients[0],
			User:      suite.users[1],
		},
	}
	for _, testRefreshToken := range testRefreshTokens {
		err = suite.db.Create(testRefreshToken).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	// Authenticate with the first access token
	now1 := time.Now()
	accessToken, err = suite.service.Authenticate("test_token_1")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "test_token_1", accessToken.Token)
	assert.EqualValues(suite.T(), suite.clients[0].ID, accessToken.ClientID.Int64)
	assert.EqualValues(suite.T(), suite.users[0].ID, accessToken.UserID.Int64)

	// First refresh token expiration date should be extended
	refreshTokens = make([]*RefreshToken, len(testRefreshTokens))
	err = suite.service.db.Where(
		"token IN ('test_token_1', 'test_token_2', 'test_token_3')",
	).Order("created_at").Find(&refreshTokens).Error
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "test_token_1", refreshTokens[0].Token)
	assert.Equal(
		suite.T(),
		now1.Unix()+int64(suite.cnf.Oauth.RefreshTokenLifetime),
		refreshTokens[0].ExpiresAt.Unix(),
	)
	assert.Equal(suite.T(), "test_token_2", refreshTokens[1].Token)
	assert.Equal(
		suite.T(),
		testRefreshTokens[1].ExpiresAt.Unix(),
		refreshTokens[1].ExpiresAt.Unix(),
	)
	assert.Equal(suite.T(), "test_token_3", refreshTokens[2].Token)
	assert.Equal(
		suite.T(),
		testRefreshTokens[2].ExpiresAt.Unix(),
		refreshTokens[2].ExpiresAt.Unix(),
	)

	// Authenticate with the second access token
	now2 := time.Now()
	accessToken, err = suite.service.Authenticate("test_token_2")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "test_token_2", accessToken.Token)
	assert.EqualValues(suite.T(), suite.clients[0].ID, accessToken.ClientID.Int64)
	assert.False(suite.T(), accessToken.UserID.Valid)

	// Second refresh token expiration date should be extended
	refreshTokens = make([]*RefreshToken, len(testRefreshTokens))
	err = suite.service.db.Where(
		"token IN ('test_token_1', 'test_token_2', 'test_token_3')",
	).Order("created_at").Find(&refreshTokens).Error
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "test_token_1", refreshTokens[0].Token)
	assert.Equal(
		suite.T(),
		now1.Unix()+int64(suite.cnf.Oauth.RefreshTokenLifetime),
		refreshTokens[0].ExpiresAt.Unix(),
	)
	assert.Equal(suite.T(), "test_token_2", refreshTokens[1].Token)
	assert.Equal(
		suite.T(),
		now2.Unix()+int64(suite.cnf.Oauth.RefreshTokenLifetime),
		refreshTokens[1].ExpiresAt.Unix(),
	)
	assert.Equal(suite.T(), "test_token_3", refreshTokens[2].Token)
	assert.Equal(
		suite.T(),
		testRefreshTokens[2].ExpiresAt.Unix(),
		refreshTokens[2].ExpiresAt.Unix(),
	)

	// Authenticate with the third access token
	now3 := time.Now()
	accessToken, err = suite.service.Authenticate("test_token_3")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "test_token_3", accessToken.Token)
	assert.EqualValues(suite.T(), suite.clients[0].ID, accessToken.ClientID.Int64)
	assert.EqualValues(suite.T(), suite.users[1].ID, accessToken.UserID.Int64)

	// First refresh token expiration date should be extended
	refreshTokens = make([]*RefreshToken, len(testRefreshTokens))
	err = suite.service.db.Where(
		"token IN ('test_token_1', 'test_token_2', 'test_token_3')",
	).Order("created_at").Find(&refreshTokens).Error
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "test_token_1", refreshTokens[0].Token)
	assert.Equal(
		suite.T(),
		now1.Unix()+int64(suite.cnf.Oauth.RefreshTokenLifetime),
		refreshTokens[0].ExpiresAt.Unix(),
	)
	assert.Equal(suite.T(), "test_token_2", refreshTokens[1].Token)
	assert.Equal(
		suite.T(),
		now2.Unix()+int64(suite.cnf.Oauth.RefreshTokenLifetime),
		refreshTokens[1].ExpiresAt.Unix(),
	)
	assert.Equal(suite.T(), "test_token_3", refreshTokens[2].Token)
	assert.Equal(
		suite.T(),
		now3.Unix()+int64(suite.cnf.Oauth.RefreshTokenLifetime),
		refreshTokens[2].ExpiresAt.Unix(),
	)
}
