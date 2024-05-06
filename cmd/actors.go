package main

import (
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/makooster/MCA/pkg/model"
	"github.com/makooster/MCA/pkg/validator"
)

func (app *application) getActorListHandler(w http.ResponseWriter, r *http.Request){
	// Embed the new Filters struct.

	var input struct {
		Fullname    string `json:"title"`
		FilmID      int    `json:"film_id"`
		model.Filters
	}

	// Initialize a new Validator instance.
	v := validator.New()
	
	// Call r.URL.Query() to get the url.Values map containing the query string data.
	qs := r.URL.Query()

	input.Fullname = app.readString(qs, "fullname", "")

	input.FilmID = app.readInt(qs, "film_id", 1, v)
	// Get the page and page_size query string values as integers. Notice that we set
	// the default page value to 1 and default page_size to 20, and that we pass the
	// validator instance as the final argument here.
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	// Extract the sort query string value, falling back to "id" if it is not provided
	// by the client (which will imply a ascending sort on movie ID).
	input.Filters.Sort = app.readString(qs, "sort", "id")
	// Add the supported sort values for this endpoint to the sort safelist.
	input.Filters.SortSafelist = []string{"id", "full_name","dorama_id", "-id", "-fullname","-dorama_id"}

	// Execute the validation checks on the Filters struct and send a response
	// containing the errors if necessary.
	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Call the GetAll() method to retrieve the movies, passing in the various filter
	// parameters.
	// Accept the metadata struct as a return value.
	
	actors, metadata, err := app.models.Doramas.GetAll(input.Fullname, input.FilmID, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Send a JSON response containing the movie data.
	// Include the metadata in the response envelope.
	
	err = app.writeJSON(w, http.StatusOK, envelope{"actors": actors, "metadata": metadata}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createActorHandler(w http.ResponseWriter, r *http.Request) {
	var input model.Actor

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.Actors.Insert(&input)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, input)
}

func (app *application) getActorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid actor ID")
		return
	}

	actor, err := app.models.Actors.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, actor)
}

func (app *application) updateActorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid actor ID")
		return
	}

	actor, err := app.models.Actors.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input model.Actor

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	actor.Name = input.Name

	err = app.models.Actors.Update(actor)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, actor)
}

func (app *application) deleteActorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid actor ID")
		return
	}

	err = app.models.Actors.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
