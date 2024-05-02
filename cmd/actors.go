package main

import (
	"net/http"
	"strconv"
	"errors"
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/makooster/MCA/pkg/model"
)

func (app *application) ActorHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.getActorsHandler(w, r)
	case http.MethodPost:
		app.createActorHandler(w, r)
	default:
		app.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// newly added getActorsHandler

func (app *application) getActorsHandler(w http.ResponseWriter, r *http.Request){
    var input model.Actor

    err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

    actor, err := app.models.Actors.Get(input.ActorId) 
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            app.respondWithError(w, http.StatusNotFound, "Actor not found")
            return
        } else {
            app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
            return
        }
    }

    app.respondWithJSON(w, http.StatusOK, actor)
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
