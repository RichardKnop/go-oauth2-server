package oauth2

import (
	"encoding/json"

	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func (suite *TestSuite) TestRegisterUsernameRequired() {
	r := test.MakeSimpleRequest(
		"POST", "http://1.2.3.4/oauth2/api/v1/users",
		map[string]interface{}{
			"username":   "",
			"password":   "test_password",
			"first_name": "John",
			"last_name":  "Doe",
		},
	)
	recorded := test.RunRequest(suite.T(), suite.API.MakeHandler(), r)

	// Status code
	assert.Equal(
		suite.T(),
		400,
		recorded.Recorder.Code, "Status code should be 400",
	)

	// Response body
	assert.Equal(
		suite.T(),
		"{\"error\":\"username required\"}",
		recorded.Recorder.Body.String(),
		"Body should be expected JSON error",
	)
}

func (suite *TestSuite) TestRegisterPasswordRequired() {
	r := test.MakeSimpleRequest(
		"POST", "http://1.2.3.4/oauth2/api/v1/users",
		map[string]interface{}{
			"username":   "test_username",
			"password":   "",
			"first_name": "John",
			"last_name":  "Doe",
		},
	)
	recorded := test.RunRequest(suite.T(), suite.API.MakeHandler(), r)

	// Status code
	assert.Equal(
		suite.T(),
		400,
		recorded.Recorder.Code, "Status code should be 400",
	)

	// Response body
	assert.Equal(
		suite.T(),
		"{\"error\":\"password required\"}",
		recorded.Recorder.Body.String(),
		"Body should be expected JSON error",
	)
}

func (suite *TestSuite) TestRegisterFirstNameRequired() {
	r := test.MakeSimpleRequest(
		"POST", "http://1.2.3.4/oauth2/api/v1/users",
		map[string]interface{}{
			"username":   "test_username",
			"password":   "test_password",
			"first_name": "",
			"last_name":  "Doe",
		},
	)
	recorded := test.RunRequest(suite.T(), suite.API.MakeHandler(), r)

	// Status code
	assert.Equal(
		suite.T(),
		400,
		recorded.Recorder.Code, "Status code should be 400",
	)

	// Response body
	assert.Equal(
		suite.T(),
		"{\"error\":\"first_name required\"}",
		recorded.Recorder.Body.String(),
		"Body should be expected JSON error",
	)
}

func (suite *TestSuite) TestRegisterLastNameNameRequired() {
	r := test.MakeSimpleRequest(
		"POST", "http://1.2.3.4/oauth2/api/v1/users",
		map[string]interface{}{
			"username":   "test_username",
			"password":   "test_password",
			"first_name": "John",
			"last_name":  "",
		},
	)
	recorded := test.RunRequest(suite.T(), suite.API.MakeHandler(), r)

	// Status code
	assert.Equal(
		suite.T(),
		400,
		recorded.Recorder.Code, "Status code should be 400",
	)

	// Response body
	assert.Equal(
		suite.T(),
		"{\"error\":\"last_name required\"}",
		recorded.Recorder.Body.String(),
		"Body should be expected JSON error",
	)
}

func (suite *TestSuite) TestRegisterUsernameAlreadyTaken() {
	r := test.MakeSimpleRequest(
		"POST", "http://1.2.3.4/oauth2/api/v1/users",
		map[string]interface{}{
			"username":   "test_USERname", // test case insensitivity of usernames
			"password":   "test_password",
			"first_name": "John",
			"last_name":  "Doe",
		},
	)
	recorded := test.RunRequest(suite.T(), suite.API.MakeHandler(), r)

	// Status code
	assert.Equal(
		suite.T(),
		400,
		recorded.Recorder.Code, "Status code should be 400",
	)

	// Response body
	assert.Equal(
		suite.T(),
		"{\"error\":\"test_USERname already taken\"}",
		recorded.Recorder.Body.String(),
		"Body should be expected JSON error",
	)
}

func (suite *TestSuite) TestRegister() {
	r := test.MakeSimpleRequest(
		"POST", "http://1.2.3.4/oauth2/api/v1/users",
		map[string]interface{}{
			"username":   "test_username_2",
			"password":   "test_password_2",
			"first_name": "John",
			"last_name":  "Doe",
		},
	)
	recorded := test.RunRequest(suite.T(), suite.API.MakeHandler(), r)

	// Status code
	assert.Equal(
		suite.T(),
		200,
		recorded.Recorder.Code, "Status code should be 200",
	)

	// Response body
	expected, _ := json.Marshal(map[string]interface{}{
		"id":         2,
		"username":   "test_username_2",
		"first_name": "John",
		"last_name":  "Doe",
	})
	assert.Equal(
		suite.T(),
		string(expected),
		recorded.Recorder.Body.String(),
		"Response body should be expected user object",
	)

	// User record was inserted
	user := User{}
	assert.Equal(
		suite.T(),
		false,
		suite.DB.Where("LOWER(username) = LOWER(?)", "test_username_2").First(&user).RecordNotFound(),
		"User should be in the database",
	)

	// Password properly hashed
	assert.Nil(
		suite.T(),
		bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("test_password_2")),
	)
}
