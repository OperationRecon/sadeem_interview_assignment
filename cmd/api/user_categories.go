package main

import (
	"net/http"
)

// Handle setting and updating category relations with a user
func (app *application) setRelationsHandler(w http.ResponseWriter, r *http.Request) {
	// strcuture input
	var input struct {
		UserID     int `json:"user_id"`
		CategoryID int `json:"category_id"`
	}

	// read input
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.UserCategories.InsertUserCategories(input.UserID, input.CategoryID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"message": "relations added"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) deleteRelationsHandler(w http.ResponseWriter, r *http.Request) {
	// strcuture input
	var input struct {
		UserID     int `json:"user_id"`
		CategoryID int `json:"category_id"`
	}

	// read input
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	app.models.UserCategories.DeleteUserCategories(input.UserID, input.CategoryID)

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "relations removed"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
