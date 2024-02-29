package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

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
	err = app.writeJSON(w, http.StatusCreated, env, nil)
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

	// Get the email for the user whose pfp is to be added
	email, err := app.readEmailParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	// use email to fetch other relevant info
	// Fetch user info from database
	user, err := app.models.Users.UserGet(email, *r)
	if err != nil {
		if err == data.ErrRecordNotFound {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	// Make sure that the sent reuqest conatins an image
	v := validator.New()
	if v.Check(validator.PermitedFileType(mimeType, "image/jpeg", "image/jpg", "image/png"),
		"image", "Improper File Type"); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Validation succesful, attempt to add image

	// Save image into server
	// get percise file type
	fType := strings.Split(mimeType, "/")

	// create thepath to be used for storing the image in the server
	fName := fmt.Sprintf("%s%d.%s", user.Name, user.ID, fType[len(fType)-1])

	// Create image filepath
	fPath := path.Join(app.config.pictureDir, fName)

	// Create file
	_, err = os.Create(fPath)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = os.WriteFile(fPath, byte, os.ModeAppend)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Save image filepath into database
	err = app.models.Users.UserUpdatePicture(fName, user.Email)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// fetch user one last time ensure updated info
	// Fetch user info from database
	user, err = app.models.Users.UserGet(user.Email, *r)
	if err != nil {
		if err == data.ErrRecordNotFound {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"Message": "Image added succesfully",
		"User":    user,
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
	user, err := app.models.Users.UserGet(email, *r)
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

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Read values to use in updating the user

	// / Create the input structure
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

	// get the email to be updated
	email, err := app.readEmailParam(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
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

	// All good? update user information
	err = app.models.Users.UserUpdate(*user, email)
	if err != nil {
		switch {
		case err == data.ErrDuplicateEmail:
			{
				app.badRequestResponse(w, r, err)
				return
			}

		case err == data.ErrRecordNotFound:
			{
				app.notFoundResponse(w, r)
				return
			}

		}
	}

	// fecth updated data
	*user, err = app.models.Users.UserGet(user.Email, *r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Compose JSON reply
	envelope := envelope{
		"message": "updated information:",
		"user":    user,
	}

	err = app.writeJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get the Email to use for user lookup
	email, err := app.readEmailParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Use email for query
	err = app.models.Users.UserDelete(email)
	if err != nil {
		if err == data.ErrRecordNotFound {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	// Infrom user of deletion
	envelope := envelope{
		"message": "User deleted successfully",
	}

	err = app.writeJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
