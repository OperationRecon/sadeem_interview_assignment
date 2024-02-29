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
	router.HandlerFunc(http.MethodGet, "/v1/users/:email", app.requireAdmin(app.requireAuthenticatedUser(app.getUserHandler)))
	router.HandlerFunc(http.MethodPut, "/v1/users/:email", app.requireAdmin(app.requireAuthenticatedUser(app.updateUserHandler)))
	router.HandlerFunc(http.MethodDelete, "/v1/users/:email", app.requireAdmin(app.requireAuthenticatedUser(app.deleteUserHandler)))
	router.HandlerFunc(http.MethodPost, "/v1/users", app.createUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/:email/pfpicture", app.requireAdmin(app.requireAuthenticatedUser(app.insertImageHandler)))
	// usercategory relations methods
	router.HandlerFunc(http.MethodDelete, "/v1/user_categories", app.requireAdmin(app.requireAuthenticatedUser(app.deleteRelationsHandler)))
	router.HandlerFunc(http.MethodPut, "/v1/user_categories", app.requireAdmin(app.requireAuthenticatedUser(app.setRelationsHandler)))
	// Category methods
	router.HandlerFunc(http.MethodPost, "/v1/categories", app.requireAdmin(app.requireAuthenticatedUser(app.createCategoryHandler)))
	router.HandlerFunc(http.MethodGet, "/v1/categories", app.requireAuthenticatedUser(app.getCategoriesHandler))
	router.HandlerFunc(http.MethodPut, "/v1/categories/:id", app.requireAdmin(app.requireAuthenticatedUser(app.updateCategoryHandler)))
	router.HandlerFunc(http.MethodDelete, "/v1/categories/:id", app.requireAdmin(app.requireAuthenticatedUser(app.deleteCategoryHandler)))

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	// Return the httprouter instance.
	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
