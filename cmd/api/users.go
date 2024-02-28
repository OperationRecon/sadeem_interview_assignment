package main

import (
	"io"
	"net/http"

	"interview_assignment.mohamednaas.net/internal/data"
	"interview_assignment.mohamednaas.net/internal/validator"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new user and add them to the database

	// Create the input structure
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Picture  string `json:"picture"`
	}

	// Read the request into the input struct
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	// Read the input into appropriate structure
	user := &data.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Picture:  input.Picture,
	}
	// Validate the clients input
	v := validator.New()

	if data.ValidateUserRegisteration(v, user); !v.Valid() {
		// Validation failed, send appropriate error messages
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Validation succesful, attempt to create user.
	// Inserting user into database
	id, err := app.models.Users.UserCreate(*user)

	if err != nil {
		if err == data.ErrDuplicateEmail {
			app.badRequestResponse(w, r, err)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	env := envelope{
		"meassage": "User Created Sucessfully",
		"id":       id,
	}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) insertImageHandler(w http.ResponseWriter, r *http.Request) {

	// Make sure to close the request file after being done with processing the image
	defer r.Body.Close()

	byte, err := io.ReadAll(r.Body)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}
	mimeType := http.DetectContentType(byte)

	// Make sure that the sent reuqest conatins an image
	v := validator.New()
	if v.Check(validator.PermitedFileType(mimeType, "image/jpeg", "image/png"),
		"image", "Improper File Type"); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Validation succesful, attempt to add image
	env := envelope{
		"Message": "Create a new image",
	}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// Handler for fecthing user information via email
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {

	email, err := app.readEmailParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	// Fetch user info from database
	user, err := app.models.Users.UserGet(email)
	if err != nil {
		if err == data.ErrRecordNotFound {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	// All good? wrap and output user info
	envelope := envelope{
		"message": "user information",
		"user":    user,
	}
	err = app.writeJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
