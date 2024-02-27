package data

import (
	"database/sql"

	"interview_assignment.mohamednaas.net/internal/validator"
)

// the user model used for connecting user info with the databse
type UserModel struct {
	DB *sql.DB
}

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Picture  string `json:"picture"`
}

// Perform checks to make sure that the registered user is valid
func ValidateUserRegisteration(v *validator.Validator, u *User) {
	v.Check(validator.NotBlank(u.Name), "name", "Name must be provided")
	v.Check(validator.NotBlank(u.Email), "email", "Email must be provided")
	v.Check(validator.NotBlank(u.Password), "password", "Password must be provided")
	v.Check(validator.Matches(u.Email, validator.EmailRX), "email", "Email must be a valid address")
	v.Check(validator.MinChars(u.Password, 8), "password", "Password must be atleast 8 characters long")

}

// Inserting a user into the database, returns newly created user's id
func (m *UserModel) UserCreate(u User) int {
	return 0
}

// Getting user info from the database
func (m *UserModel) UserGet(id int) User {
	return User{}
}

// Updating user info

// Adding a profile picture to the server and database
func (m *UserModel) UserUpdatePicture(u User) error {
	return nil
}

// Updateing other user info using a JSON request
func (m *UserModel) UserUpdate(u User) error {
	return nil
}

// Deleting a User by id
func (m *UserModel) UserDelete(id int) error {
	return nil
}
