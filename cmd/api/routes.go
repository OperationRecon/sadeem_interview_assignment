package main

import (
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	// Initialize a new httprouter router instance.
	router := httprouter.New()

	// Register the relevant methods, URL patterns and handler functions for our
	// endpoints using the HandlerFunc() method.

	// handle serving the static files
	router.ServeFiles("/static/*filepath", http.Dir(os.Getenv("sainpr_pfp_dir")))

	// Set custom handlers for aftermentioned routes
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	// User Methods
	router.HandlerFunc(http.MethodGet, "/v1/users/:email", app.requireAuthenticatedUser(app.getUserHandler))
	router.HandlerFunc(http.MethodPut, "/v1/users/:email", app.requireAuthenticatedUser(app.updateUserHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/users/:email", app.requireAuthenticatedUser(app.deleteUserHandler))
	router.HandlerFunc(http.MethodPost, "/v1/users", app.createUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/:email/pfpicture", app.requireAuthenticatedUser(app.insertImageHandler))
	// Category methods
	router.HandlerFunc(http.MethodPost, "/v1/categories", app.requireAuthenticatedUser(app.createCategoryHandler))
	router.HandlerFunc(http.MethodGet, "/v1/categories", app.requireAuthenticatedUser(app.getCategoriesHandler))
	router.HandlerFunc(http.MethodPut, "/v1/categories/:id", app.requireAuthenticatedUser(app.updateCategoryHandler))

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	// Return the httprouter instance.
	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
