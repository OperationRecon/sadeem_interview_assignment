package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"interview_assignment.mohamednaas.net/userdata"
)

func (app *application) routes() *httprouter.Router {
	// Initialize a new httprouter router instance.
	router := httprouter.New()

	// Register the relevant methods, URL patterns and handler functions for our
	// endpoints using the HandlerFunc() method.

	// create file server to handle serveing out of ui/static/
	fileServer := http.FileServer(http.FS(userdata.Files))

	// handle serving the static files
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	// Set custom handlers for aftermentioned routes
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/users/:email", app.getUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users", app.createUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/:email/pfpicture", app.insertImageHandler)

	// Return the httprouter instance.
	return router
}
