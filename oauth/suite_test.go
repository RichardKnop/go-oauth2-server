package oauth_test

import (
	"log"
	"os"
	"testing"

	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/RichardKnop/go-oauth2-server/database"
	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
)

var testDbPath = "/tmp/oauth_testdb.sqlite"

var testFixtures = []string{
	"./oauth/fixtures/scopes.yml",
	"./oauth/fixtures/test_clients.yml",
	"./oauth/fixtures/test_users.yml",
}

// db migrations needed for tests
var testMigrations = []func(*gorm.DB) error{
	oauth.MigrateAll,
}

func init() {
	if err := os.Chdir("../"); err != nil {
		log.Fatal(err)
	}
}

// OauthTestSuite needs to be exported so the tests run
type OauthTestSuite struct {
	suite.Suite
	cnf     *config.Config
	db      *gorm.DB
	service *oauth.Service
	clients []*oauth.Client
	users   []*oauth.User
	router  *mux.Router
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (suite *OauthTestSuite) SetupSuite() {

	// Initialise the config
	suite.cnf = config.NewConfig(false, false)

	// Create the test database
	db, err := database.CreateTestDatabase(testDbPath, testMigrations, testFixtures)
	if err != nil {
		log.Fatal(err)
	}
	suite.db = db

	// Fetch test client
	suite.clients = make([]*oauth.Client, 0)
	if err := suite.db.Order("id").Find(&suite.clients).Error; err != nil {
		log.Fatal(err)
	}

	// Fetch test users
	suite.users = make([]*oauth.User, 0)
	if err := suite.db.Order("id").Find(&suite.users).Error; err != nil {
		log.Fatal(err)
	}

	// Initialise the service
	suite.service = oauth.NewService(suite.cnf, suite.db)

	// Register routes
	suite.router = mux.NewRouter()
	oauth.RegisterRoutes(suite.router, suite.service)
}

// The TearDownSuite method will be run by testify once, at the very
// end of the testing suite, after all tests have been run.
func (suite *OauthTestSuite) TearDownSuite() {
	//
}

// The SetupTest method will be run before every test in the suite.
func (suite *OauthTestSuite) SetupTest() {
	//
}

// The TearDownTest method will be run after every test in the suite.
func (suite *OauthTestSuite) TearDownTest() {
	// Scopes are static, populated from fixtures,
	// so there is no need to clear them after running a test
	suite.db.Unscoped().Delete(new(oauth.AuthorizationCode))
	suite.db.Unscoped().Delete(new(oauth.RefreshToken))
	suite.db.Unscoped().Delete(new(oauth.AccessToken))
	suite.db.Unscoped().Not("id", []int64{1, 2}).Delete(new(oauth.User))
	suite.db.Unscoped().Not("id", []int64{1, 2, 3}).Delete(new(oauth.Client))
}

// TestOauthTestSuite ...
// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestOauthTestSuite(t *testing.T) {
	suite.Run(t, new(OauthTestSuite))
}
