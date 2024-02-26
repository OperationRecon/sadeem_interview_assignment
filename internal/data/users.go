package data

import "interview_assignment.mohamednaas.net/internal/validator"

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Picture  string `json:"picture"`
}

// Perform checks to make sure that the registered user is valid
func ValidateUserRegister(v *validator.Validator, u *User) {
	v.Check(validator.NotBlank(u.Name), "name", "Name must be provided")
	v.Check(validator.NotBlank(u.Email), "email", "Email must be provided")
	v.Check(validator.NotBlank(u.Password), "password", "Password must be provided")
	v.Check(validator.Matches(u.Email, validator.EmailRX), "email", "Email must be a valid address")
	v.Check(validator.MinChars(u.Password, 8), "password", "Password must be atleast 8 characters long")

}
