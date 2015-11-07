package oauth

import (
	"log"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/RichardKnop/go-oauth2-server/api"
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/RichardKnop/go-oauth2-server/migrate"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
	// sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

// OAuth2TestSuite ...
type OAuth2TestSuite struct {
	suite.Suite
	DB  *gorm.DB
	API *rest.Api
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (suite *OAuth2TestSuite) SetupSuite() {
	// Init in-memory test database
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	suite.DB = &db
	migrate.Bootstrap(&db)
	MigrateAll(&db)

	// Init API app
	suite.API = api.NewAPI(
		api.DevelopmentStack,
		NewRoutes(config.NewConfig(), suite.DB),
	)
}

// The TearDownSuite method will be run by testify once, at the very
// end of the testing suite, after all tests have been run.
func (suite *OAuth2TestSuite) TearDownSuite() {
	//
}

// The SetupTest method will be run before every test in the suite.
func (suite *OAuth2TestSuite) SetupTest() {
	// Insert test client
	clientSecretHash, err := bcrypt.GenerateFromPassword([]byte("test_client_secret"), 3)
	if err != nil {
		log.Fatal(err)
	}
	if err := suite.DB.Create(&Client{
		ID:       1,
		ClientID: "test_client_id",
		Password: string(clientSecretHash),
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Insert test user
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("test_password"), 3)
	if err != nil {
		log.Fatal(err)
	}
	if err := suite.DB.Create(&User{
		ID:       1,
		Username: "test_username",
		Password: string(passwordHash),
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Insert test scopes
	if err := suite.DB.Create(&Scope{
		Scope:     "foo",
		IsDefault: true,
	}).Error; err != nil {
		log.Fatal(err)
	}
	if err := suite.DB.Create(&Scope{
		Scope:     "bar",
		IsDefault: true,
	}).Error; err != nil {
		log.Fatal(err)
	}
	if err := suite.DB.Create(&Scope{
		Scope:     "qux",
		IsDefault: false,
	}).Error; err != nil {
		log.Fatal(err)
	}
}

// The TearDownTest method will be run after every test in the suite.
func (suite *OAuth2TestSuite) TearDownTest() {
	// Empty all the tables
	suite.DB.Delete(AccessToken{})
	suite.DB.Delete(RefreshToken{})
	suite.DB.Delete(AuthCode{})
	suite.DB.Delete(Scope{})
	suite.DB.Delete(User{})
	suite.DB.Delete(Client{})
}

// TestTestSuite ...
// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestTestSuite(t *testing.T) {
	suite.Run(t, new(OAuth2TestSuite))
}
