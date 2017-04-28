package session_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/RichardKnop/go-oauth2-server/session"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// SessionTestSuite needs to be exported so the tests run
type SessionTestSuite struct {
	suite.Suite
	cnf     *config.Config
	service *session.Service
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (suite *SessionTestSuite) SetupSuite() {
	suite.cnf = config.NewConfig(false, false, "etcd")

	// Overwrite internal vars so we don't affect existing session
	session.StorageSessionName = "test_session"
	session.UserSessionKey = "test_user"

	// Initialise the service
	r, err := http.NewRequest("GET", "http://1.2.3.4/foo/bar", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	w := httptest.NewRecorder()

	suite.service = session.NewService(suite.cnf, sessions.NewCookieStore([]byte(suite.cnf.Session.Secret)))
	suite.service.SetSessionService(r, w)
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

// TesSessionTestSuite ...
// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TesSessionTestSuite(t *testing.T) {
	suite.Run(t, new(SessionTestSuite))
}
