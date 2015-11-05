package oauth2

import (
	"net/url"
	"testing"

	"github.com/RichardKnop/go-microservice-example/api"
	"github.com/RichardKnop/go-microservice-example/config"
	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestRegisterUser() {
	api := api.NewAPI(NewRoutes(config.NewConfig(), suite.DB))
	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/api/v1/users",
		url.Values{"username": {"testusername"}, "password": {"testpassword"}})
	recorded := test.RunRequest(suite.T(), api.MakeHandler(), r)

	assert.Equal(suite.T(), 200, recorded.Recorder.Code, "Status code should be 200")
	expectedBody := "{\n  \"id\": 1,\n  \"username\": \"testusername\"\n}"
	assert.Equal(suite.T(), expectedBody, recorded.Recorder.Body.String(),
		"Body should be expected JSON error")
}

// TestTestSuite
// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
