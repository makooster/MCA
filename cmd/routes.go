package main

import (
	"net/http"
	"github.com/gorilla/mux"
)

// "github.com/julienschmidt/httprouter"

func (app *application) routes()  http.Handler {
	// Initialize a new httprouter router instance.
	// router := httprouter.New()
	router := mux.NewRouter()
	// Convert the notFoundResponse() helper to a http.Handler using the
	// http.HandlerFunc() adapter, and then set it as the custom error handler for 404
	// Not Found responses.
	router.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)
	
	// Likewise, convert the methodNotAllowedResponse() helper to a http.Handler and set
	// it as the custom error handler for 405 Method Not Allowed responses.
	router.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	// Register the relevant methods, URL patterns and handler functions for our
	// endpoints using the HandlerFunc() method. Note that http.MethodGet and
	// http.MethodPost are constants which equate to the strings "GET" and "POST"
	// respectively.
	router.HandleFunc("/app/check", app.healthcheckHandler).Methods("GET")

	router.HandleFunc("/app/doramas", app.requirePermission("movies:read", app.getDoramaListHandler)).Methods("GET")
	router.HandleFunc("/app/doramas", app.requirePermission("movies:write", app.createDoramaHandler)).Methods("POST")
	router.HandleFunc("/app/doramas/{id:[0-9]+}", app.requirePermission("movies:read", app.getDoramaHandler)).Methods("GET")
	router.HandleFunc("/app/doramas/{id:[0-9]+}", app.requirePermission("movies:write", app.updateDoramaHandler)).Methods("PUT")
	router.HandleFunc("/app/doramas/{id:[0-9]+}", app.requirePermission("movies:write", app.deleteDoramaHandler)).Methods("DELETE")

	router.HandleFunc("/app/actors", app.requirePermission("movies:read", app.getActorListHandler)).Methods("GET")
	router.HandleFunc("/app/actors", app.requirePermission("movies:write", app.createActorHandler)).Methods("POST")
	router.HandleFunc("/app/actors/{id:[0-9]+}", app.requirePermission("movies:read", app.getActorHandler)).Methods("GET")
	router.HandleFunc("/app/actors/{id:[0-9]+}", app.requirePermission("movies:write", app.updateActorHandler)).Methods("PUT")
	router.HandleFunc("/app/actors/{id:[0-9]+}", app.requirePermission("movies:write", app.deleteActorHandler)).Methods("DELETE")

	router.HandleFunc("/app/genres", app.requirePermission("movies:read", app.getGenresListHandler)).Methods("GET")
	router.HandleFunc("/app/genres", app.requirePermission("movies:write", app.createGenreHandler)).Methods("POST")
	router.HandleFunc("/app/genres/{id:[0-9]+}", app.requirePermission("movies:read", app.getGenreHandler)).Methods("GET")
	router.HandleFunc("/app/genres/{id:[0-9]+}", app.requirePermission("movies:write", app.updateGenreHandler)).Methods("PUT")
	router.HandleFunc("/app/genres/{id:[0-9]+}", app.requirePermission("movies:write", app.deleteGenreHandler)).Methods("DELETE")

	router.HandleFunc("/app/users", app.registerUserHandler).Methods("POST")
	router.HandleFunc("/app/users/activated", app.activateUserHandler).Methods("PUT")
	router.HandleFunc("/app/tokens/login", app.createAuthenticationTokenHandler).Methods("POST")

	// return router
	return app.authenticate(router)

}
