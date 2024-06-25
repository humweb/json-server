package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/humweb/json-server/internal/storage"
	"net/http"
	"net/http/httptest"
)

var (
	mockServer *httptest.Server

	testResourceKeys    = []string{"apples", "oranges"}
	testData            = make(storage.Database)
	testResourceStorage = make(map[string]*storage.Mock)
)

//func TestMain(m *testing.M) {
//
//	resourceKeys, err := testGenerateData()
//	if err != nil {
//		panic(err)
//	}
//
//	resourceStorage, err := testCreateResourceStorage(resourceKeys)
//	if err != nil {
//		panic(err)
//	}
//
//	router := handler.Setup(resourceStorage, false, false)
//
//	mockServer = httptest.NewServer(router)
//	defer mockServer.Close()
//
//	os.Exit(m.Run())
//}

func testCreateResourceStorage(resourceKeys []string) (map[string]storage.Storage, error) {
	resourceStorage := make(map[string]storage.Storage)

	for _, resourceKey := range resourceKeys {
		storageSvc, err := storage.NewMock(testData, resourceKey)
		if err != nil {
			return nil, errors.New("failed to initialize resources")
		}

		resourceStorage[resourceKey] = storageSvc
	}

	for key, storageSvc := range resourceStorage {
		testResourceStorage[key] = storageSvc.(*storage.Mock)
	}

	return resourceStorage, nil
}

func testListResourcesByKey(key string) ([]storage.Resource, error) {
	url := fmt.Sprintf("%s/%s", mockServer.URL, key)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var body []storage.Resource
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body, nil
}
