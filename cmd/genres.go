package main

import (
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/makooster/MCA/pkg/model"
	"github.com/makooster/MCA/pkg/validator"
)

func (app *application) getGenresListHandler(w http.ResponseWriter, r *http.Request){
	// Embed the new Filters struct.

	var input struct {
		GenreName      string `json:"genre_name"`
		GenreID        int    `json:"genre_id"`
		model.Filters
	}

	// Initialize a new Validator instance.
	v := validator.New()
	
	// Call r.URL.Query() to get the url.Values map containing the query string data.
	qs := r.URL.Query()

	input.GenreName = app.readString(qs, "genre_name", "")

	input.GenreID = app.readInt(qs, "genre_id", 1, v)
	// Get the page and page_size query string values as integers. Notice that we set
	// the default page value to 1 and default page_size to 20, and that we pass the
	// validator instance as the final argument here.
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	// Extract the sort query string value, falling back to "id" if it is not provided
	// by the client (which will imply a ascending sort on movie ID).
	input.Filters.Sort = app.readString(qs, "sort", "genre_id")
	// Add the supported sort values for this endpoint to the sort safelist.
	input.Filters.SortSafelist = []string{"genre_id", "genre_name","-genre_id", "-genre_name"}

	// Execute the validation checks on the Filters struct and send a response
	// containing the errors if necessary.
	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Call the GetAll() method to retrieve the movies, passing in the various filter
	// parameters.
	// Accept the metadata struct as a return value.
	
	genres, metadata, err := app.models.Actors.GetAll(input.GenreName, input.GenreID, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Send a JSON response containing the movie data.
	// Include the metadata in the response envelope.
	
	err = app.writeJSON(w, http.StatusOK, envelope{"genres": genres, "metadata": metadata}, nil)
	
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getGenreHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *application) createGenreHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *application) updateGenreHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]
	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid genre ID")
		return
	}

	genre, err := app.models.Genres.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input model.Genre

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
 
	genre.GenreName = input.GenreName
	
	
	err = app.models.Genres.Update(genre)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, genre)
}

func (app *application) deleteGenreHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid genre ID")
		return
	}

	err = app.models.Actors.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
