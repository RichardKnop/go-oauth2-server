package oauth

import (
	"log"
	"testing"

	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/RichardKnop/go-oauth2-server/migrations"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
	// sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

// OauthTestSuite needs to be exported so the tests run
type OauthTestSuite struct {
	suite.Suite
	cnf     *config.Config
	db      *gorm.DB
	service *Service
	client  *Client
	user    *User
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (suite *OauthTestSuite) SetupSuite() {
	suite.cnf = config.NewConfig()

	// Init in-memory test database
	inMemoryDB, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	suite.db = &inMemoryDB

	// Run all migrations
	migrations.Bootstrap(suite.db)
	MigrateAll(suite.db)

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
	// Insert a test client
	client, err := suite.service.CreateClient(
		"test_client",             // client id
		"test_secret",             // client secret
		"https://www.example.com", // redirect URI
	)
	if err != nil {
		log.Fatal(err)
	}
	suite.client = client

	// Insert a test user
	user, err := suite.service.CreateUser(
		"test@username", // username
		"test_password", // password
	)
	if err != nil {
		log.Fatal(err)
	}
	suite.user = user

	// Insert test scopes
	testScopes := map[string]bool{
		"read":       true,
		"read_write": false,
	}
	for scope, isDefault := range testScopes {
		if err := suite.db.Create(&Scope{
			Scope:     scope,
			IsDefault: isDefault,
		}).Error; err != nil {
			log.Fatal(err)
		}
	}
}

// The TearDownTest method will be run after every test in the suite.
func (suite *OauthTestSuite) TearDownTest() {
	// Empty all the tables
	suite.db.Unscoped().Delete(AuthorizationCode{})
	suite.db.Unscoped().Delete(RefreshToken{})
	suite.db.Unscoped().Delete(AccessToken{})
	suite.db.Unscoped().Delete(Scope{})
	suite.db.Unscoped().Delete(User{})
	suite.db.Unscoped().Delete(Client{})
}

// TestOauthTestSuite ...
// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestOauthTestSuite(t *testing.T) {
	suite.Run(t, new(OauthTestSuite))
}
