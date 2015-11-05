package oauth2

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func (suite *TestSuite) TestRegisterUsernameRequired() {
	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/api/v1/users", nil)
	r.PostForm = url.Values{}
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
		"{\"Error\":\"username required\"}",
		recorded.Recorder.Body.String(),
		"Body should be expected JSON error",
	)
}

func (suite *TestSuite) TestRegisterPasswordRequired() {
	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/api/v1/users", nil)
	r.PostForm = url.Values{
		"username": {"testusername"},
	}
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
		"{\"Error\":\"password required\"}",
		recorded.Recorder.Body.String(),
		"Body should be expected JSON error",
	)
}

func (suite *TestSuite) TestRegisterFirstNameRequired() {
	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/api/v1/users", nil)
	r.PostForm = url.Values{
		"username": {"testusername"},
		"password": {"testpassword"},
	}
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
		"{\"Error\":\"first_name required\"}",
		recorded.Recorder.Body.String(),
		"Body should be expected JSON error",
	)
}

func (suite *TestSuite) TestRegisterLastNameNameRequired() {
	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/api/v1/users", nil)
	r.PostForm = url.Values{
		"username":   {"testusername"},
		"password":   {"testpassword"},
		"first_name": {"John"},
	}
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
		"{\"Error\":\"last_name required\"}",
		recorded.Recorder.Body.String(),
		"Body should be expected JSON error",
	)
}

func (suite *TestSuite) TestRegister() {
	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/api/v1/users", nil)
	r.PostForm = url.Values{
		"username":   {"testusername"},
		"password":   {"testpassword"},
		"first_name": {"John"},
		"last_name":  {"Doe"},
	}
	recorded := test.RunRequest(suite.T(), suite.API.MakeHandler(), r)

	// Status code
	assert.Equal(
		suite.T(),
		200,
		recorded.Recorder.Code, "Status code should be 200",
	)

	// Response body
	expected, _ := json.Marshal(map[string]interface{}{
		"id":         1,
		"username":   "testusername",
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
		suite.DB.Where("LOWER(username) = LOWER(?)", "testusername").First(&user).RecordNotFound(),
		"User should be in the database",
	)

	// Password properly hashed
	assert.Nil(
		suite.T(),
		bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("testpassword")),
	)
}

func (suite *TestSuite) TestRegisterUsernameAlreadyTaken() {
	if err := suite.DB.Create(&User{
		Username:  "testUSERname",
		Password:  "doesn't matter",
		FirstName: "doesn't matter",
		LastName:  "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}

	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/api/v1/users", nil)
	r.PostForm = url.Values{
		"username":   {"testusername"},
		"password":   {"testpassword"},
		"first_name": {"John"},
		"last_name":  {"Doe"},
	}
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
		"{\"Error\":\"testusername already taken\"}",
		recorded.Recorder.Body.String(),
		"Body should be expected JSON error",
	)
}
