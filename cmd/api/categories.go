package main

import (
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

	err = app.writeJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Sends all categories to user, sorted by id
func (app *application) getCategoriesHandler(w http.ResponseWriter, r *http.Request) {

	// get all categories from DB
	categories, err := app.models.Categories.CategoriesGet()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// prepare Json response
	envelope := envelope{
		"message":    "all categories:",
		"categories": categories,
	}

	err = app.writeJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
