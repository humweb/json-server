// Package handler contains the full set of handler functions and routes
// supported by the web api.
package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/humweb/json-server/internal/storage"
	"github.com/rs/cors"
	"net/http"
)

// Setup API handler based on provided resources.
func Setup(resourceStorage storage.Map, shouldLog, allowAll bool) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	if shouldLog {
		router.Use(middleware.Logger)
	}

	if allowAll {
		router.Use(cors.AllowAll().Handler)
	}

	// For each resource create the appropriate endpoint handlers.
	for resourceKey, storageService := range resourceStorage {
		handlers := NewAPI(storageService)

		// Common endpoint to retrieve db contents.
		if resourceKey == "db" {
			router.Get("/db", handlers.DB)
			continue
		}

		// Register all default endpoint handlers for resource.
		router.Options(fmt.Sprintf("/%s", resourceKey), handlers.Options)
		router.Options(fmt.Sprintf("/%s/{id}", resourceKey), handlers.Options)
		router.Get(fmt.Sprintf("/%s", resourceKey), handlers.Find)
		router.Get(fmt.Sprintf("/%s/{id}", resourceKey), handlers.FindBy)
		router.Post(fmt.Sprintf("/%s", resourceKey), handlers.Create)
		router.Put(fmt.Sprintf("/%s/{id}", resourceKey), handlers.Replace)
		router.Patch(fmt.Sprintf("/%s/{id}", resourceKey), handlers.Update)
		router.Delete(fmt.Sprintf("/%s/{id}", resourceKey), handlers.Delete)
	}

	// Render a home page with useful info.
	router.Get("/", HomePage(resourceStorage))

	return router
}
