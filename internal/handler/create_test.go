package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/humweb/json-server/internal/handler"
	"github.com/humweb/json-server/internal/storage"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type HandlerCreateTestSuite struct {
	suite.Suite
	store    *storage.File
	handlers *handler.API
	db       storage.Database
}
type TestResponse struct {
	Id     int    `json:"id"`
	Field1 string `json:"field_1"`
	Field2 string `json:"field_2"`
}

func (suite *HandlerCreateTestSuite) SetupSuite() {
}
func (suite *HandlerCreateTestSuite) SetupTest() {
	f, _ := testGenerateStorageFile()
	suite.store, _ = storage.NewFile(f.Name(), "users")
	suite.handlers = handler.NewAPI(suite.store)

}

func (suite *HandlerCreateTestSuite) TestCreate() {

	body := storage.Resource{
		"id":       "11",
		"name":     "new-field_1",
		"username": "new-field_2",
	}
	w, r := mockRequest("POST", "/users", body, nil)

	suite.handlers.Create(w, r)

	res, err := suite.store.FindById("11")

	suite.Nil(err)
	suite.Equal(res, body)
}
func (suite *HandlerCreateTestSuite) TestWithoutId() {

	body := storage.Resource{
		"field_1": "new-field_1",
		"field_2": "new-field_2",
	}
	w, r := mockRequest("POST", "/users", body, nil)

	suite.handlers.Create(w, r)

	res, err := suite.store.FindById("10")
	var newBody TestResponse
	//spew.Dump("BODY", json.Unmarshal(r.Response.Body, newBody))
	suite.Nil(err)
	suite.Equal(body, res)
}

func TestInertiaHttpSuite(t *testing.T) {
	suite.Run(t, new(HandlerCreateTestSuite))
}

type Headers = map[string]string

func mockRequest(method string, target string, body any, headers Headers) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	var r *http.Request

	b, _ := json.Marshal(body)
	r = httptest.NewRequest(method, target, bytes.NewBuffer(b))

	for key, val := range headers {
		r.Header.Set(key, val)
	}

	return w, r
}

func testGenerateStorageFile() (*os.File, error) {
	f, err := os.CreateTemp(".", "")
	if err != nil {
		return nil, err
	}

	contentBytes, err := testGenerateData()
	if err != nil {
		return nil, err
	}

	if err = os.WriteFile(f.Name(), contentBytes, 0644); err != nil {
		return nil, err
	}

	return f, nil
}

func testGenerateData() ([]byte, error) {

	testData := make(storage.Database)
	keys := []string{"users"}

	for _, key := range keys {
		resources := []storage.Resource{}
		for idx := 0; idx < 11; idx++ {
			newResource := storage.Resource{
				"id":       strconv.Itoa(idx),
				"name":     fmt.Sprintf("name_-%s-%d", key, idx),
				"username": fmt.Sprintf("username_-%s-%d", key, idx),
			}

			resources = append(resources, newResource)
		}

		testData[key] = resources
	}

	return json.Marshal(testData)
}

func testResetData(filename string) error {
	contentBytes, err := json.Marshal(testData)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filename, contentBytes, 0644); err != nil {
		return err
	}

	return nil
}
