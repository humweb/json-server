// Package storage defines an interface which describes storage CRUD functionality.
package storage

import (
	"errors"
	"math/rand"
	"strconv"
)

var (
	// ErrResourceNotFound returns an error when a requested resource not found in storage.
	ErrResourceNotFound = errors.New("resource not found")
	// ErrResourceAlreadyExists returns an error when a resource already exists in storage.
	ErrResourceAlreadyExists = errors.New("resource already exists")
	// ErrBadRequest returns an error when an unexpected request been processed.
	ErrBadRequest = errors.New("bad request")
	// ErrInternalServerError returns an error when an unexpected error occurs.
	ErrInternalServerError = errors.New("internal Server Error")
)

// Resource represents the structure of a singe resource in storage.
type Resource map[string]interface{}

// Database represents the structure of the storage contents.
type Database map[string][]Resource

type Map map[string]Storage

// Storage interface to handle storage operations.
type Storage interface {
	Find() ([]Resource, error)
	FindById(string) (Resource, error)
	Create(Resource) (Resource, error)
	Replace(string, Resource) (Resource, error)
	Update(string, Resource) (Resource, error)
	Delete(string) error
	DB() (Database, error)
}

// generateNewId and validate that is unique across provided data.
func generateNewId(data []Resource) string {

	existingIds := make(map[string]bool)
	for _, d := range data {
		existingIds[d["id"].(string)] = true
	}

	for {
		newId := strconv.Itoa(rand.Intn(1000))

		if !existingIds[newId] {
			return newId
		}
	}
}

// checkResourceKeyExists in the file data.
func checkResourceKeyExists(database Database, key string) error {
	if _, ok := database[key]; !ok {
		return ErrResourceNotFound
	}

	return nil
}
