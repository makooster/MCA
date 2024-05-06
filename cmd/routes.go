package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes()  http.Handler {
	// Initialize a new httprouter router instance.
	router := httprouter.New()

	// Convert the notFoundResponse() helper to a http.Handler using the
	// http.HandlerFunc() adapter, and then set it as the custom error handler for 404
	// Not Found responses.
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	
	// Likewise, convert the methodNotAllowedResponse() helper to a http.Handler and set
	// it as the custom error handler for 405 Method Not Allowed responses.
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Register the relevant methods, URL patterns and handler functions for our
	// endpoints using the HandlerFunc() method. Note that http.MethodGet and
	// http.MethodPost are constants which equate to the strings "GET" and "POST"
	// respectively.
	router.HandlerFunc(http.MethodGet, "/app/check", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/app/doramas", app.requirePermission("movies:read", app.getDoramaListHandler))
	router.HandlerFunc(http.MethodPost, "/app/doramas", app.requirePermission("movies:write", app.createDoramaHandler))
	router.HandlerFunc(http.MethodGet, "/app/doramas/:id", app.requirePermission("movies:read", app.getDoramaHandler))
	router.HandlerFunc(http.MethodPut, "/app/doramas/:id", app.requirePermission("movies:write", app.updateDoramaHandler))
	router.HandlerFunc(http.MethodDelete, "/app/doramas/:id", app.requirePermission("movies:write", app.deleteDoramaHandler))
	
	router.HandlerFunc(http.MethodGet, "/app/actors", app.requirePermission("movies:read", app.getActorListHandler))
	router.HandlerFunc(http.MethodPost, "/app/actors", app.requirePermission("movies:write", app.createActorHandler))
	router.HandlerFunc(http.MethodGet, "/app/actors/:id", app.requirePermission("movies:read", app.getActorHandler))
	router.HandlerFunc(http.MethodPut, "/app/actors/:id", app.requirePermission("movies:write", app.updateActorHandler))
	router.HandlerFunc(http.MethodDelete,"/app/actors/:id", app.requirePermission("movies:write", app.deleteActorHandler))

	router.HandlerFunc(http.MethodPost, "/app/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/app/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/app/tokens/login", app.createAuthenticationTokenHandler)



	// Return the httprouter instance.
	// return router
	return app.authenticate(router)

}
