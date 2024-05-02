package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) UserHandler(w http.ResponseWriter, r *http.Request) {
	// Handle user page
	app.respondWithJSON(w, http.StatusOK, map[string]string{"message": "User page"})
}

func (app *application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	app.respondWithJSON(w, http.StatusOK, map[string]string{"message": "Welcome to the home page"})
}

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJSON(w, code, map[string]string{"error": message})
}

func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}