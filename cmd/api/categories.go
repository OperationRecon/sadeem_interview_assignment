package main

import (
	"errors"
	"net/http"

	"interview_assignment.mohamednaas.net/internal/data"
	"interview_assignment.mohamednaas.net/internal/validator"
)

// Handle the creation of a new category
func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// read input parameters into a nice struct
	var input struct {
		Name string `json:"name"`
	}

	// Read input from json request
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// place input into category struct
	category := &data.Category{
		Name: input.Name,
	}

	// Validate input
	v := validator.New()

	if data.ValidateCategoryInsertion(v, category); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
	}

	// Valid input paramters, insert category into database
	category.ID, err = app.models.Categories.CategoryCreate(*category)
	if err != nil {
		if err == data.ErrDuplicateCategoryName {
			app.badRequestResponse(w, r, err)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	// All good? prepare Json response
	envelope := envelope{
		"message":  "Category created successfully",
		"category": category,
	}

	err = app.writeJSON(w, http.StatusCreated, envelope, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Sends all categories to user, sorted by id
func (app *application) getCategoriesHandler(w http.ResponseWriter, r *http.Request) {

	// prepare response envelope

	var (
		err        error
		categories []*data.Category
	)
	// Check if user is admin
	user := app.contextGetUser(r)

	if app.models.Users.IsAdmin(user.ID) {
		// get all categories from DB

		categories, err = app.models.Categories.CategoriesGet()
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

	} else {
		categories, err = app.models.UserCategories.UserCategoriesGet(user.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

	}
	err = app.writeJSON(w, http.StatusOK, envelope{"categories": categories}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Fetch the existing record from the database, sending a 404 Not Found
	// response to the client if we couldn't find a matching record.
	category, err := app.models.Categories.CategoryGet(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Declare an input struct to hold the expected data from the client.
	var input struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	}
	// Read the JSON request body data into the input struct.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Copy the values from the request body to the appropriate fields of the category
	// record.
	category.Name = input.Name
	category.ID = id
	// Validate the updated record, sending the client a 422 Unprocessable Entity
	// response if any checks fail.
	v := validator.New()
	if data.ValidateCategoryInsertion(v, &category); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Categories.CategoryUpdate(category)
	if err != nil {
		if err == data.ErrDuplicateCategoryName {
			app.badRequestResponse(w, r, err)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	// Write the updated  record in a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"category": category}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Get the id delete the specfied category
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// ensure category to be deleted exists
	_, err = app.models.Categories.CategoryGet(id)
	if err != nil {
		if err == data.ErrRecordNotFound {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	app.models.Categories.CategoryDelete(id)

	// Delete sucessful, write response
	app.writeJSON(w, http.StatusOK, envelope{"message": "category deleted successfully"}, nil)

}
