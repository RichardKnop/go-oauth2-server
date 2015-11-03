package service

import (
	"testing"

	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/stretchr/testify/mock"
)

type MockedDB struct {
	mock.Mock
}

func (m *MockedDB) DoSomething(number int) (bool, error) {

	args := m.Called(number)
	return args.Bool(0), args.Error(1)

}

func TestNewDatabasePostgres(t *testing.T) {
	r := test.MakeSimpleRequest("POST", "http://127.0.0.1/api/v1/users", nil)
	recorded := test.RunRequest(t, RegisterUser, r)

	//fmt.Printf("%d - %s", w.Code, w.Body.String())
}
