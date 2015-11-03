package service

import (
	"testing"

	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/mock"
)

// type MockedDB struct {
// 	mock.Mock
// }
//
// func (m *MockedDB) DoSomething(number int) (bool, error) {
//
// 	args := m.Called(number)
// 	return args.Bool(0), args.Error(1)
//
// }

func TestRegisterUserBadPayload(t *testing.T) {
	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/api/v1/users", nil)
	recorded := test.RunRequest(t, NewAPI().MakeHandler(), r)

	assert.Equal(t, 400, recorded.Recorder.Code, "Status code should be 400")
	expectedBody := "{\n  \"Error\": \"Decode JSON error\"\n}"
	assert.Equal(t, expectedBody, recorded.Recorder.Body.String(),
		"Body should be expected JSON error")
}
