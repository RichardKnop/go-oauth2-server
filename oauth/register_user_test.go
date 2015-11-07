package oauth

// import (
// 	"encoding/json"
//
// 	"github.com/ant0ine/go-json-rest/rest/test"
// 	"github.com/stretchr/testify/assert"
// 	"golang.org/x/crypto/bcrypt"
// )
//
// func (suite *OAuth2TestSuite) TestRegisterUsernameAlreadyTaken() {
// 	r := test.MakeSimpleRequest(
// 		"POST", "http://1.2.3.4/oauth2/api/v1/users",
// 		map[string]interface{}{
// 			"username": "test_USERname", // test case insensitivity of usernames
// 			"password": "test_password",
// 		},
// 	)
// 	recorded := test.RunRequest(suite.T(), suite.API.MakeHandler(), r)
//
// 	// Status code
// 	assert.Equal(
// 		suite.T(),
// 		400,
// 		recorded.Recorder.Code, "Status code should be 400",
// 	)
//
// 	// Response body
// 	assert.Equal(
// 		suite.T(),
// 		"{\"error\":\"test_USERname already taken\"}",
// 		recorded.Recorder.Body.String(),
// 		"Body should be expected JSON error",
// 	)
// }
//
// func (suite *OAuth2TestSuite) TestRegister() {
// 	r := test.MakeSimpleRequest(
// 		"POST", "http://1.2.3.4/oauth2/api/v1/users",
// 		map[string]interface{}{
// 			"username": "test_username_2",
// 			"password": "test_password_2",
// 		},
// 	)
// 	recorded := test.RunRequest(suite.T(), suite.API.MakeHandler(), r)
//
// 	// Status code
// 	assert.Equal(
// 		suite.T(),
// 		200,
// 		recorded.Recorder.Code, "Status code should be 200",
// 	)
//
// 	// User record was inserted
// 	user := User{}
// 	assert.Equal(
// 		suite.T(),
// 		false,
// 		suite.DB.Where("LOWER(username) = LOWER(?)", "test_username_2").First(&user).RecordNotFound(),
// 		"User should be in the database",
// 	)
//
// 	// Response body
// 	expected, _ := json.Marshal(map[string]interface{}{
// 		"id":       user.ID,
// 		"username": "test_username_2",
// 	})
// 	assert.Equal(
// 		suite.T(),
// 		string(expected),
// 		recorded.Recorder.Body.String(),
// 		"Response body should be expected user object",
// 	)
//
// 	// Password properly hashed
// 	assert.Nil(
// 		suite.T(),
// 		bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("test_password_2")),
// 	)
// }
