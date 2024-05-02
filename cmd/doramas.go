package main

import (
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/makooster/MCA/pkg/model"
	"github.com/makooster/MCA/pkg/validator"
)

func (app *application) getDoramaListHandler(w http.ResponseWriter, r *http.Request){
	// Embed the new Filters struct.

	var input struct {
		Title    string `json:"title"`
		ReleaseYear int `json:"release_year"`
		Duration    int `json:"duration"`
		model.Filters
	}

	// Initialize a new Validator instance.
	v := validator.New()
	
	// Call r.URL.Query() to get the url.Values map containing the query string data.
	qs := r.URL.Query()
	
	// Use our helpers to extract the title and genres query string values, falling back
	// to defaults of an empty string and an empty slice respectively if they are not
	// provided by the client.
	input.Title = app.readString(qs, "title", "")

	input.ReleaseYear = app.readInt(qs, "release_year", 1, v)

	// input.Duration = app.readInt(qs, "duration", 1, v)
	

	// Get the page and page_size query string values as integers. Notice that we set
	// the default page value to 1 and default page_size to 20, and that we pass the
	// validator instance as the final argument here.
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	// Extract the sort query string value, falling back to "id" if it is not provided
	// by the client (which will imply a ascending sort on movie ID).
	input.Filters.Sort = app.readString(qs, "sort", "dorama_id")
	// Add the supported sort values for this endpoint to the sort safelist.
	input.Filters.SortSafelist = []string{"dorama_id", "title","release_year", "-dorama_id", "-title","-release_year"}

	// Execute the validation checks on the Filters struct and send a response
	// containing the errors if necessary.
	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Call the GetAll() method to retrieve the movies, passing in the various filter
	// parameters.
	// Accept the metadata struct as a return value.
	
	doramas, metadata, err := app.models.Doramas.GetAll(input.Title, input.ReleaseYear,input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Send a JSON response containing the movie data.
	// Include the metadata in the response envelope.
	
	err = app.writeJSON(w, http.StatusOK, envelope{"doramas": doramas, "metadata": metadata}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// func (app *application) getDoramasHandler(w http.ResponseWriter, r *http.Request){
// 	var input model.Dorama

// 	err := app.readJSON(w, r, &input)
// 	if err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	dorama, err := app.models.Doramas.Get(input.DoramaId) 
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			app.respondWithError(w, http.StatusNotFound, "Dorama not found")
// 			return
// 		} else {
// 			app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
// 			return
// 		}
// 	}

// 	app.respondWithJSON(w, http.StatusOK, dorama)
// }

func (app *application) getDoramaHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid drama ID")
		return
	}

	drama, err := app.models.Doramas.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, drama)
}



func (app *application) createDoramaHandler(w http.ResponseWriter, r *http.Request) {
	var input model.Dorama

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	
	err = app.models.Doramas.Insert(&input)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, input)
}

func (app *application) updateDoramaHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid drama ID")
		return
	}

	dorama, err := app.models.Doramas.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input model.Dorama

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	dorama.Title = input.Title
	dorama.Description = input.Description
	dorama.ReleaseYear = input.ReleaseYear
	dorama.MainActors = input.MainActors

	err = app.models.Doramas.Update(dorama)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, dorama)
}

func (app *application) deleteDoramaHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid drama ID")
		return
	}

	err = app.models.Doramas.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}