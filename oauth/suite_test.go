package oauth

import (
	"log"
	"testing"

	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/RichardKnop/go-oauth2-server/database"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
)

var testDbPath = "/tmp/oauth_testdb.sqlite"

var testFixtures = []string{
	"fixtures/scopes.yml",
	"fixtures/test_clients.yml",
	"fixtures/test_users.yml",
}

// db migrations needed for tests
var testMigrations = []func(*gorm.DB) error{
	MigrateAll,
}

// OauthTestSuite needs to be exported so the tests run
type OauthTestSuite struct {
	suite.Suite
	cnf     *config.Config
	db      *gorm.DB
	service *Service
	clients []*Client
	users   []*User
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
	suite.clients = make([]*Client, 0)
	if err := suite.db.Order("id").Find(&suite.clients).Error; err != nil {
		log.Fatal(err)
	}

	// Fetch test users
	suite.users = make([]*User, 0)
	if err := suite.db.Order("id").Find(&suite.users).Error; err != nil {
		log.Fatal(err)
	}

	// Initialise the service
	suite.service = NewService(suite.cnf, suite.db)
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
	suite.db.Unscoped().Delete(new(AuthorizationCode))
	suite.db.Unscoped().Delete(new(RefreshToken))
	suite.db.Unscoped().Delete(new(AccessToken))
	suite.db.Unscoped().Not("id", []int64{1, 2}).Delete(new(User))
	suite.db.Unscoped().Not("id", []int64{1, 2}).Delete(new(Client))
}

// TestOauthTestSuite ...
// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestOauthTestSuite(t *testing.T) {
	suite.Run(t, new(OauthTestSuite))
}
