package oauth2

// import (
// 	"testing"
//
// 	"github.com/RichardKnop/go-microservice-example/api"
// 	"github.com/RichardKnop/go-microservice-example/config"
// 	"github.com/RichardKnop/go-microservice-example/utils"
// 	"github.com/ant0ine/go-json-rest/rest/test"
// 	"github.com/stretchr/testify/assert"
// )
//
// func TestRegisterUserBadPayload(t *testing.T) {
// 	db := utils.TestSetUp()
//
// 	api := api.NewAPI(NewRoutes(config.NewConfig(), db))
// 	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/api/v1/users", nil)
// 	recorded := test.RunRequest(t, api.MakeHandler(), r)
//
// 	assert.Equal(t, 400, recorded.Recorder.Code, "Status code should be 400")
// 	expectedBody := "{\n  \"Error\": \"Decode JSON error\"\n}"
// 	assert.Equal(t, expectedBody, recorded.Recorder.Body.String(),
// 		"Body should be expected JSON error")
//
// 	utils.TestTearDown(db)
// }
