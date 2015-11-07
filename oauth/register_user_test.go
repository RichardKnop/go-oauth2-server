package oauth

// import (
// 	"encoding/json"
//
// 	"github.com/ant0ine/go-json-rest/rest/test"
// 	"github.com/stretchr/testify/assert"
// 	"golang.org/x/crypto/bcrypt"
// )
//
// func (suite *OauthTestSuite) TestRegisterUsernameAlreadyTaken() {
//  // Make a request
// 	r := test.MakeSimpleRequest(
// 		"POST", "http://1.2.3.4/oauth2/api/v1/users",
// 		map[string]interface{}{
// 			"username": "test_USERname", // test case insensitivity of usernames
// 			"password": "test_password",
// 		},
// 	)
// 	recorded := test.RunRequest(suite.T(), suite.API.MakeHandler(), r)
//
// 	// Check the status code
// 	assert.Equal(
// 		suite.T(),
// 		400,
// 		recorded.Recorder.Code, "Status code should be 400",
// 	)
//
// 	// Check the response body
// 	assert.Equal(
// 		suite.T(),
// 		"{\"error\":\"test_USERname already taken\"}",
// 		recorded.Recorder.Body.String(),
// 		"Body should be expected JSON error",
// 	)
// }
//
// func (suite *OauthTestSuite) TestRegister() {
//  // Make a request
// 	r := test.MakeSimpleRequest(
// 		"POST", "http://1.2.3.4/oauth2/api/v1/users",
// 		map[string]interface{}{
// 			"username": "test_username_2",
// 			"password": "test_password_2",
// 		},
// 	)
// 	recorded := test.RunRequest(suite.T(), suite.API.MakeHandler(), r)
//
// 	// Check the status code
// 	assert.Equal(
// 		suite.T(),
// 		200,
// 		recorded.Recorder.Code, "Status code should be 200",
// 	)
//
// 	// Check the user record was inserted
// 	user := User{}
// 	assert.Equal(
// 		suite.T(),
// 		false,
// 		suite.DB.Where("LOWER(username) = LOWER(?)", "test_username_2").First(&user).RecordNotFound(),
// 		"User should be in the database",
// 	)
//
// 	// Check the response body
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
// 	// Check the password was properly hashed
// 	assert.Nil(
// 		suite.T(),
// 		bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("test_password_2")),
// 	)
// }
