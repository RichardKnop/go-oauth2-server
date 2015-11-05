package oauth2

import (
	"log"
	"testing"

	"github.com/RichardKnop/go-microservice-example/api"
	"github.com/RichardKnop/go-microservice-example/config"
	"github.com/RichardKnop/go-microservice-example/migrate"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
	// sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

// TestSuite ...
type TestSuite struct {
	suite.Suite
	DB  *gorm.DB
	API *rest.Api
}

// SetupTest creates in-memory test database and starts app
func (suite *TestSuite) SetupTest() {
	if suite.DB == nil {
		db, err := gorm.Open("sqlite3", ":memory:")
		if err != nil {
			log.Fatal(err)
		}
		suite.DB = &db
		migrate.Bootstrap(&db)
		MigrateAll(&db)
	}

	if suite.API == nil {
		stack := []rest.Middleware{
			&rest.AccessLogApacheMiddleware{},
			&rest.TimerMiddleware{},
			&rest.RecorderMiddleware{},
			&rest.PoweredByMiddleware{},
			&rest.RecoverMiddleware{
				EnableResponseStackTrace: true,
			},
		}
		suite.API = api.NewAPI(
			stack,
			NewRoutes(config.NewConfig(), suite.DB),
		)
	}
}

// TearDown truncates all tables
func (suite *TestSuite) TearDown() {
	suite.DB.Exec("DELETE FROM SELECT name FROM sqlite_master WHERE type IS 'table'")
}

// TestTestSuite ...
// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
