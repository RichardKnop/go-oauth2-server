package session_test

import (
	"github.com/RichardKnop/go-oauth2-server/session"
	"github.com/stretchr/testify/assert"
)

func (suite *SessionTestSuite) TestService() {
	var (
		userSession *session.UserSession
		err         error
	)

	// No public methods should work before StartSession has been called
	userSession, err = suite.service.GetUserSession()
	assert.Nil(suite.T(), userSession)
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), session.ErrSessonNotStarted, err)
	}

	// Call the StartSession method so internal session object gets set
	err = suite.service.StartSession()
	assert.Nil(suite.T(), err)

	// Let's clear the user session first
	err = suite.service.ClearUserSession()
	assert.Nil(suite.T(), err)

	// Since the user session has not been set yet, this should return error
	userSession, err = suite.service.GetUserSession()
	assert.Nil(suite.T(), userSession)
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "User session type assertion error", err.Error())
	}

	// Let's set the user session now
	suite.service.SetUserSession(&session.UserSession{
		ClientID:     "test_client",
		Username:     "test@username",
		AccessToken:  "test_access_token",
		RefreshToken: "test_refresh_token",
	})

	// User session is set, this should return it
	userSession, err = suite.service.GetUserSession()
	assert.Nil(suite.T(), err)
	if assert.NotNil(suite.T(), userSession) {
		assert.Equal(suite.T(), "test_client", userSession.ClientID)
		assert.Equal(suite.T(), "test@username", userSession.Username)
		assert.Equal(suite.T(), "test_access_token", userSession.AccessToken)
		assert.Equal(suite.T(), "test_refresh_token", userSession.RefreshToken)
	}

	// Let's clear the user session now
	err = suite.service.ClearUserSession()
	assert.Nil(suite.T(), err)

	// Since the user session has been cleared, this should return error
	userSession, err = suite.service.GetUserSession()
	assert.Nil(suite.T(), userSession)
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "User session type assertion error", err.Error())
	}
}
