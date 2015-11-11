package oauth

import (
	"log"
	"testing"

	"golang.org/x/crypto/bcrypt"

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
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	suite.db = &db

	// Run all migrations
	migrations.Bootstrap(suite.db)
	MigrateAll(suite.db)

	suite.service = &Service{cnf: suite.cnf, db: suite.db}
}

// The TearDownSuite method will be run by testify once, at the very
// end of the testing suite, after all tests have been run.
func (suite *OauthTestSuite) TearDownSuite() {
	//
}

// The SetupTest method will be run before every test in the suite.
func (suite *OauthTestSuite) SetupTest() {
	// Insert a test client
	clientSecretHash, err := bcrypt.GenerateFromPassword([]byte("test_secret"), 3)
	if err != nil {
		log.Fatal(err)
	}
	suite.client = &Client{
		ID:       1,
		ClientID: "test_client",
		Secret:   string(clientSecretHash),
	}
	if err := suite.db.Create(suite.client).Error; err != nil {
		log.Fatal(err)
	}

	// Insert a test user
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("test_password"), 3)
	if err != nil {
		log.Fatal(err)
	}
	suite.user = &User{
		ID:       1,
		Username: "test_username",
		Password: string(passwordHash),
	}
	if err := suite.db.Create(suite.user).Error; err != nil {
		log.Fatal(err)
	}

	// Insert test scopes
	if err := suite.db.Create(&Scope{
		Scope:     "foo",
		IsDefault: true,
	}).Error; err != nil {
		log.Fatal(err)
	}
	if err := suite.db.Create(&Scope{
		Scope:     "bar",
		IsDefault: true,
	}).Error; err != nil {
		log.Fatal(err)
	}
	if err := suite.db.Create(&Scope{
		Scope:     "qux",
		IsDefault: false,
	}).Error; err != nil {
		log.Fatal(err)
	}
}

// The TearDownTest method will be run after every test in the suite.
func (suite *OauthTestSuite) TearDownTest() {
	// Empty all the tables
	suite.db.Delete(AccessToken{})
	suite.db.Delete(RefreshToken{})
	suite.db.Delete(AuthorizationCode{})
	suite.db.Delete(Scope{})
	suite.db.Delete(User{})
	suite.db.Delete(Client{})
}

// TestOauthTestSuite ...
// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestOauthTestSuite(t *testing.T) {
	suite.Run(t, new(OauthTestSuite))
}
