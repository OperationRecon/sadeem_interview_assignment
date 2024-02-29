package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"interview_assignment.mohamednaas.net/internal/validator"
)

var (
	ErrDuplicateCategoryName = errors.New("duplicate category name")
)

// the Category model used for connecting category info with the databse
type CategoryModel struct {
	DB *sql.DB
}

type Category struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

func ValidateCategoryInsertion(v *validator.Validator, c *Category) {
	v.Check(validator.NotBlank(c.Name), "name", "Name must be provided")
}

// Categoty insertion
func (m *CategoryModel) CategoryCreate(c Category) (int, error) {
	// prepare query
	q := "INSERT INTO categories (name) VALUES ($1) RETURNING id"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// If the table already contains a record with this email address, then when we try
	// to perform the insert there will be a violation of the UNIQUE "users_email_key"
	// constraint that we set up in the previous chapter. We check for this error
	// specifically, and return custom ErrDuplicateEmail error instead.
	err := m.DB.QueryRowContext(ctx, q, c.Name).Scan(&c.ID)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "categories_name_key"`:
			return 0, ErrDuplicateCategoryName
		default:
			return 0, err
		}
	}
	return c.ID, nil
}

// fetch all categories
func (m *CategoryModel) CategoriesGet() ([]*Category, error) {
	// Prepare query
	q := `SELECT * FROM categories ORDER BY id`

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use QueryContext() to execute the query. This returns a sql.Rows resultset
	// containing the result.
	rows, err := m.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	// Importantly, defer a call to rows.Close() to ensure that the resultset is closed
	// before GetAll() returns.
	defer rows.Close()
	// Initialize an empty slice to hold the category data.
	categories := []*Category{}
	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty Category struct to hold the data for an individual category.
		var c Category
		// Scan the values from the row into the Category struct.
		err := rows.Scan(
			&c.ID,
			&c.Name,
		)
		if err != nil {
			return nil, err
		}
		// Add the Category struct to the slice.
		categories = append(categories, &c)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK, then return the slice of categories.
	return categories, nil
}

// Category Update
func (m *CategoryModel) CategoryUpdate(c *Category) error {
	return nil
}

// Category Delete by name
func (m *CategoryModel) CategoryDelete(name string) error {
	return nil
}
