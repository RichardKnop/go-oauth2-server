package oauth_test

import (
	"time"

	"github.com/RichardKnop/go-oauth2-server/models"
	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/RichardKnop/go-oauth2-server/session"
	"github.com/RichardKnop/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestAuthenticate() {
	var (
		accessToken *models.OauthAccessToken
		err         error
	)

	// Insert some test access tokens
	testAccessTokens := []*models.OauthAccessToken{
		// Expired access token
		{
			MyGormModel: models.MyGormModel{
				ID:        uuid.New(),
				CreatedAt: time.Now().UTC(),
			},
			Token:     "test_expired_token",
			ExpiresAt: time.Now().UTC().Add(-10 * time.Second),
			Client:    suite.clients[0],
			User:      suite.users[0],
		},
		// Access token without a user
		{
			MyGormModel: models.MyGormModel{
				ID:        uuid.New(),
				CreatedAt: time.Now().UTC(),
			},
			Token:     "test_client_token",
			ExpiresAt: time.Now().UTC().Add(+10 * time.Second),
			Client:    suite.clients[0],
		},
		// Access token with a user
		{
			MyGormModel: models.MyGormModel{
				ID:        uuid.New(),
				CreatedAt: time.Now().UTC(),
			},
			Token:     "test_user_token",
			ExpiresAt: time.Now().UTC().Add(+10 * time.Second),
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
		assert.Equal(suite.T(), oauth.ErrAccessTokenNotFound, err)
	}

	// Test passing a bogus token
	accessToken, err = suite.service.Authenticate("bogus")

	// Access token should be nil
	assert.Nil(suite.T(), accessToken)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), oauth.ErrAccessTokenNotFound, err)
	}

	// Test passing an expired token
	accessToken, err = suite.service.Authenticate("test_expired_token")

	// Access token should be nil
	assert.Nil(suite.T(), accessToken)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), oauth.ErrAccessTokenExpired, err)
	}

	// Test passing a valid client token
	accessToken, err = suite.service.Authenticate("test_client_token")

	// Correct access token should be returned
	if assert.NotNil(suite.T(), accessToken) {
		assert.Equal(suite.T(), "test_client_token", accessToken.Token)
		assert.EqualValues(suite.T(), suite.clients[0].ID, accessToken.ClientID.String)
		assert.False(suite.T(), accessToken.UserID.Valid)
	}

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Test passing a valid user token
	accessToken, err = suite.service.Authenticate("test_user_token")

	// Correct access token should be returned
	if assert.NotNil(suite.T(), accessToken) {
		assert.Equal(suite.T(), "test_user_token", accessToken.Token)
		assert.EqualValues(suite.T(), suite.clients[0].ID, accessToken.ClientID.String)
		assert.EqualValues(suite.T(), suite.users[0].ID, accessToken.UserID.String)
	}

	// Error should be nil
	assert.Nil(suite.T(), err)
}

func (suite *OauthTestSuite) TestAuthenticateRollingRefreshToken() {
	var (
		testAccessTokens  []*models.OauthAccessToken
		testRefreshTokens []*models.OauthRefreshToken
		accessToken       *models.OauthAccessToken
		err               error
		refreshTokens     []*models.OauthRefreshToken
	)

	// Insert some test access tokens
	testAccessTokens = []*models.OauthAccessToken{
		{
			MyGormModel: models.MyGormModel{
				ID:        uuid.New(),
				CreatedAt: time.Now().UTC(),
			},
			Token:     "test_token_1",
			ExpiresAt: time.Now().UTC().Add(+10 * time.Second),
			Client:    suite.clients[0],
			User:      suite.users[0],
		},
		{
			MyGormModel: models.MyGormModel{
				ID:        uuid.New(),
				CreatedAt: time.Now().UTC(),
			},
			Token:     "test_token_2",
			ExpiresAt: time.Now().UTC().Add(+10 * time.Second),
			Client:    suite.clients[0],
		},
		{
			MyGormModel: models.MyGormModel{
				ID:        uuid.New(),
				CreatedAt: time.Now().UTC(),
			},
			Token:     "test_token_3",
			ExpiresAt: time.Now().UTC().Add(+10 * time.Second),
			Client:    suite.clients[0],
			User:      suite.users[1],
		},
	}
	for _, testAccessToken := range testAccessTokens {
		err = suite.db.Create(testAccessToken).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	// Insert some test access tokens
	testRefreshTokens = []*models.OauthRefreshToken{
		{
			MyGormModel: models.MyGormModel{
				ID:        uuid.New(),
				CreatedAt: time.Now().UTC(),
			},
			Token:     "test_token_1",
			ExpiresAt: time.Now().UTC().Add(+10 * time.Second),
			Client:    suite.clients[0],
			User:      suite.users[0],
		},
		{
			MyGormModel: models.MyGormModel{
				ID:        uuid.New(),
				CreatedAt: time.Now().UTC(),
			},
			Token:     "test_token_2",
			ExpiresAt: time.Now().UTC().Add(+10 * time.Second),
			Client:    suite.clients[0],
		},
		{
			MyGormModel: models.MyGormModel{
				ID:        uuid.New(),
				CreatedAt: time.Now().UTC(),
			},
			Token:     "test_token_3",
			ExpiresAt: time.Now().UTC().Add(+10 * time.Second),
			Client:    suite.clients[0],
			User:      suite.users[1],
		},
	}
	for _, testRefreshToken := range testRefreshTokens {
		err = suite.db.Create(testRefreshToken).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	// Authenticate with the first access token
	now1 := time.Now().UTC()
	gorm.NowFunc = func() time.Time {
		return now1
	}
	accessToken, err = suite.service.Authenticate("test_token_1")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "test_token_1", accessToken.Token)
	assert.EqualValues(suite.T(), suite.clients[0].ID, accessToken.ClientID.String)
	assert.EqualValues(suite.T(), suite.users[0].ID, accessToken.UserID.String)

	// First refresh token expiration date should be extended
	refreshTokens = make([]*models.OauthRefreshToken, len(testRefreshTokens))
	err = suite.db.Where(
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
	now2 := time.Now().UTC()
	gorm.NowFunc = func() time.Time {
		return now2
	}
	accessToken, err = suite.service.Authenticate("test_token_2")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "test_token_2", accessToken.Token)
	assert.EqualValues(suite.T(), suite.clients[0].ID, accessToken.ClientID.String)
	assert.False(suite.T(), accessToken.UserID.Valid)

	// Second refresh token expiration date should be extended
	refreshTokens = make([]*models.OauthRefreshToken, len(testRefreshTokens))
	err = suite.db.Where(
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
	now3 := time.Now().UTC()
	gorm.NowFunc = func() time.Time {
		return now3
	}
	accessToken, err = suite.service.Authenticate("test_token_3")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "test_token_3", accessToken.Token)
	assert.EqualValues(suite.T(), suite.clients[0].ID, accessToken.ClientID.String)
	assert.EqualValues(suite.T(), suite.users[1].ID, accessToken.UserID.String)

	// First refresh token expiration date should be extended
	refreshTokens = make([]*models.OauthRefreshToken, len(testRefreshTokens))
	err = suite.db.Where(
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

func (suite *OauthTestSuite) TestClearUserTokens() {
	var (
		testAccessTokens  []*models.OauthAccessToken
		testRefreshTokens []*models.OauthRefreshToken
		err               error
		testUserSession   *session.UserSession
	)

	// Insert some test access tokens
	testAccessTokens = []*models.OauthAccessToken{
		{
			MyGormModel: models.MyGormModel{
				ID:        uuid.New(),
				CreatedAt: time.Now().UTC(),
			},
			Token:     "test_token_1",
			ExpiresAt: time.Now().UTC().Add(+10 * time.Second),
			Client:    suite.clients[0],
			User:      suite.users[0],
		},
		{
			MyGormModel: models.MyGormModel{
				ID:        uuid.New(),
				CreatedAt: time.Now().UTC(),
			},
			Token:     "test_token_2",
			ExpiresAt: time.Now().UTC().Add(+10 * time.Second),
			Client:    suite.clients[1],
			User:      suite.users[0],
		},
		{
			MyGormModel: models.MyGormModel{
				ID:        uuid.New(),
				CreatedAt: time.Now().UTC(),
			},
			Token:     "test_token_3",
			ExpiresAt: time.Now().UTC().Add(+10 * time.Second),
			Client:    suite.clients[0],
			User:      suite.users[1],
		},
	}
	for _, testAccessToken := range testAccessTokens {
		err = suite.db.Create(testAccessToken).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	// Insert some test access tokens
	testRefreshTokens = []*models.OauthRefreshToken{
		{
			MyGormModel: models.MyGormModel{
				ID:        uuid.New(),
				CreatedAt: time.Now().UTC(),
			},
			Token:     "test_token_1",
			ExpiresAt: time.Now().UTC().Add(+10 * time.Second),
			Client:    suite.clients[0],
			User:      suite.users[0],
		},
		{
			MyGormModel: models.MyGormModel{
				ID:        uuid.New(),
				CreatedAt: time.Now().UTC(),
			},
			Token:     "test_token_2",
			ExpiresAt: time.Now().UTC().Add(+10 * time.Second),
			Client:    suite.clients[1],
			User:      suite.users[0],
		},
		{
			MyGormModel: models.MyGormModel{
				ID:        uuid.New(),
				CreatedAt: time.Now().UTC(),
			},
			Token:     "test_token_3",
			ExpiresAt: time.Now().UTC().Add(+10 * time.Second),
			Client:    suite.clients[0],
			User:      suite.users[1],
		},
	}
	for _, testRefreshToken := range testRefreshTokens {
		err = suite.db.Create(testRefreshToken).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	testUserSession = &session.UserSession{
		ClientID:     suite.clients[0].Key,
		Username:     suite.users[0].Username,
		AccessToken:  "test_token_1",
		RefreshToken: "test_token_1",
	}

	// Remove test_token_1 from accress and refresh token tables
	suite.service.ClearUserTokens(testUserSession)

	// Assert that the refresh token was removed
	found := !models.OauthRefreshTokenPreload(suite.db).Where("token = ?", testUserSession.RefreshToken).First(&models.OauthRefreshToken{}).RecordNotFound()
	assert.Equal(suite.T(), false, found)

	// Assert that the access token was removed
	found = !models.OauthAccessTokenPreload(suite.db).Where("token = ?", testUserSession.AccessToken).First(&models.OauthAccessToken{}).RecordNotFound()
	assert.Equal(suite.T(), false, found)

	// Assert that the other two tokens are still there
	// Refresh tokens
	found = !models.OauthRefreshTokenPreload(suite.db).Where("token = ?", "test_token_2").First(&models.OauthRefreshToken{}).RecordNotFound()
	assert.Equal(suite.T(), true, found)
	found = !models.OauthRefreshTokenPreload(suite.db).Where("token = ?", "test_token_3").First(&models.OauthRefreshToken{}).RecordNotFound()
	assert.Equal(suite.T(), true, found)

	// Access tokens
	found = !models.OauthAccessTokenPreload(suite.db).Where("token = ?", "test_token_2").First(&models.OauthAccessToken{}).RecordNotFound()
	assert.Equal(suite.T(), true, found)
	found = !models.OauthAccessTokenPreload(suite.db).Where("token = ?", "test_token_3").First(&models.OauthAccessToken{}).RecordNotFound()
	assert.Equal(suite.T(), true, found)

}
