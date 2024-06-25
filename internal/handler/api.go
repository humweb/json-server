package handler

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/humweb/json-server/internal/storage"
	"github.com/humweb/json-server/internal/web"
	"net/http"
)

type API struct {
	Storage storage.Storage
}

func NewAPI(storage storage.Storage) *API {
	return &API{Storage: storage}
}

func (api *API) Create(w http.ResponseWriter, r *http.Request) {
	// Read and decode request body.
	var newResource storage.Resource
	if err := json.NewDecoder(r.Body).Decode(&newResource); err != nil {
		web.Error(w, http.StatusBadRequest, storage.ErrBadRequest.Error())
		return
	}

	// Check if request body is empty, or contains only id.
	if _, ok := newResource["id"]; len(newResource) == 0 || (len(newResource) == 1 && ok) {
		web.Error(w, http.StatusBadRequest, storage.ErrBadRequest.Error())
		return
	}
	// Create the new resource.
	data, err := api.Storage.Create(newResource)

	if err != nil {
		// Already exists with the requested id.
		if errors.Is(err, storage.ErrResourceAlreadyExists) {
			web.Error(w, http.StatusConflict, err.Error())
			return
		}

		web.Error(w, http.StatusInternalServerError, storage.ErrInternalServerError.Error())
		return
	}

	web.Success(w, http.StatusCreated, data)
}

func (api *API) Delete(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	// Delete resource.
	if err := api.Storage.Delete(id); err != nil {
		// Resource not found.
		if errors.Is(err, storage.ErrResourceNotFound) {
			web.Error(w, http.StatusNotFound, err.Error())
			return
		}

		web.Error(w, http.StatusInternalServerError, storage.ErrInternalServerError.Error())
		return
	}

	web.Success(w, http.StatusOK, nil)
}

func (api *API) Find(w http.ResponseWriter, r *http.Request) {
	// Find all resources.
	data, err := api.Storage.Find()
	if err != nil {
		web.Error(w, http.StatusInternalServerError, storage.ErrInternalServerError.Error())
		return
	}

	web.Success(w, http.StatusOK, data)
}

func (api *API) FindBy(w http.ResponseWriter, r *http.Request) {
	// Read request path parameter id.
	id := chi.URLParam(r, "id")

	// Find the resource with the requested id.
	data, err := api.Storage.FindById(id)
	if err != nil {
		// Resource not found.
		if errors.Is(err, storage.ErrResourceNotFound) {
			web.Error(w, http.StatusNotFound, err.Error())
			return
		}

		web.Error(w, http.StatusInternalServerError, storage.ErrInternalServerError.Error())
		return
	}

	web.Success(w, http.StatusOK, data)
}

func (api *API) Replace(w http.ResponseWriter, r *http.Request) {
	// Read request path parameter id.
	id := chi.URLParam(r, "id")

	// Read and decode request body.
	var newResource storage.Resource
	if err := json.NewDecoder(r.Body).Decode(&newResource); err != nil {
		web.Error(w, http.StatusBadRequest, storage.ErrBadRequest.Error())
		return
	}

	// Check if request body is empty, or contains only id.
	if _, ok := newResource["id"]; len(newResource) == 0 || (len(newResource) == 1 && ok) {
		web.Error(w, http.StatusBadRequest, storage.ErrBadRequest.Error())
		return
	}

	// Replace the resource.
	data, err := api.Storage.Replace(id, newResource)
	if err != nil {
		// Resource not found.
		if errors.Is(err, storage.ErrResourceNotFound) {
			web.Error(w, http.StatusNotFound, err.Error())
			return
		}

		web.Error(w, http.StatusInternalServerError, storage.ErrInternalServerError.Error())
		return
	}

	web.Success(w, http.StatusOK, data)
}
func (api *API) Update(w http.ResponseWriter, r *http.Request) {
	// Read request path parameter id.
	id := chi.URLParam(r, "id")

	// Read and decode request body.
	var newResource storage.Resource
	if err := json.NewDecoder(r.Body).Decode(&newResource); err != nil {
		web.Error(w, http.StatusBadRequest, storage.ErrBadRequest.Error())
		return
	}

	// Check if request body is empty, or contains only id.
	if _, ok := newResource["id"]; len(newResource) == 0 || (len(newResource) == 1 && ok) {
		web.Error(w, http.StatusBadRequest, storage.ErrBadRequest.Error())
		return
	}

	// Update the resource.
	data, err := api.Storage.Update(id, newResource)
	if err != nil {
		// Resource not found.
		if errors.Is(err, storage.ErrResourceNotFound) {
			web.Error(w, http.StatusNotFound, err.Error())
			return
		}

		web.Error(w, http.StatusInternalServerError, storage.ErrInternalServerError.Error())
		return
	}

	web.Success(w, http.StatusOK, data)
}

func (api *API) Options(w http.ResponseWriter, r *http.Request) {
	return
}

func (api *API) DB(w http.ResponseWriter, r *http.Request) {
	data, err := api.Storage.DB()
	if err != nil {
		web.Error(w, http.StatusInternalServerError, storage.ErrInternalServerError.Error())
		return
	}

	web.Success(w, http.StatusOK, data)
}
