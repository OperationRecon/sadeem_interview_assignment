package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error.
var (
	ErrRecordNotFound = errors.New("record not found")
)

// A model struct to wrap around all the other models
type Models struct {
	Users      UserModel
	Categories CategoryModel
}

// For ease of use, we also add a New() method which returns a Models struct
func NewModels(db *sql.DB) Models {
	return Models{
		Users:      UserModel{DB: db},
		Categories: CategoryModel{DB: db},
	}
}
