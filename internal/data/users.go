package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"interview_assignment.mohamednaas.net/internal/validator"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

// the user model used for connecting user info with the databse
type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID       int    `json:"-"`
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
func (m *UserModel) UserCreate(u User) (int, error) {
	// Define query used
	q := `INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id`

	// Generate password hash to insert into db
	pHashed, err := Set(u.Password)
	if err != nil {
		return 0, err
	}

	args := []any{u.Name, u.Email, pHashed}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// If the table already contains a record with this email address, then when we try
	// to perform the insert there will be a violation of the UNIQUE "users_email_key"
	// constraint that we set up in the previous chapter. We check for this error
	// specifically, and return custom ErrDuplicateEmail error instead.
	err = m.DB.QueryRowContext(ctx, q, args...).Scan(&u.ID)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return 0, ErrDuplicateEmail
		default:
			return 0, err
		}
	}

	return u.ID, nil
}

// Getting user info from the database
func (m *UserModel) UserGet(email string, r http.Request) (User, error) {
	user := User{}
	// prepare query
	q := `SELECT id, name, email, pfp_filepath FROM users WHERE email = $1`

	// excecute query
	err := m.DB.QueryRow(q, email).Scan(&user.ID, &user.Name, &user.Email, &user.Picture)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrRecordNotFound
		}
		return user, err
	}

	// Make nice URl to find image in
	user.Picture = fmt.Sprintf("http://%s/static/profile_pictures/%s", r.Host, user.Picture)

	// all good? retirn user data
	return user, nil
}

// Updating user info

// Adding a profile picture to the server and database
func (m *UserModel) UserUpdatePicture(picture, email string) error {
	// Prepare Query statment
	q := `UPDATE users
		SET pfp_filepath = $1
		WHERE email = $2
		RETURNING id`

	args := []any{picture, email}

	_ = m.DB.QueryRow(q, args...)

	// insert was sucessful, carry on.
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

// The Set() method calculates the bcrypt hash of a plaintext password
func Set(plaintextPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return nil, err
	}

	return hash, err
}

// The Matches() method checks whether the provided plaintext password matches the
// hashed password
func Matches(plaintextPassword string, hash []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
