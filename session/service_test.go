package session

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// SessionTestSuite needs to be exported so the tests run
type SessionTestSuite struct {
	suite.Suite
	cnf     *config.Config
	service *Service
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (suite *SessionTestSuite) SetupSuite() {
	suite.cnf = config.NewConfig(false, false)

	// Overwrite internal vars so we don't affect existing session
	storageSessionName = "test_session"
	userSessionKey = "test_user"

	// Initialise the service
	r, err := http.NewRequest("GET", "http://1.2.3.4/foo/bar", nil)
	if err != nil {
		log.Fatal(err)
	}
	w := httptest.NewRecorder()
	suite.service = NewService(suite.cnf, r, w)
}

// The TearDownSuite method will be run by testify once, at the very
// end of the testing suite, after all tests have been run.
func (suite *SessionTestSuite) TearDownSuite() {
	//
}

// The SetupTest method will be run before every test in the suite.
func (suite *SessionTestSuite) SetupTest() {
	//
}

// The TearDownTest method will be run after every test in the suite.
func (suite *SessionTestSuite) TearDownTest() {
	//
}

func (suite *SessionTestSuite) TestService() {
	// No public methods should work before StartSession has been called
	userSession, err := suite.service.GetUserSession()
	assert.Nil(suite.T(), userSession)
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errSessonNotStarted.Error(), err.Error())
	}

	// Call the StartSession method so internal session object gets set
	if err := suite.service.StartSession(); err != nil {
		log.Fatal(err)
	}

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
	suite.service.SetUserSession(&UserSession{
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

// TesSessionTestSuite ...
// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TesSessionTestSuite(t *testing.T) {
	suite.Run(t, new(SessionTestSuite))
}
