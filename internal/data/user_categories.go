package data

import (
	"context"
	"database/sql"
	"time"
)

type UserCategoriesModel struct {
	DB *sql.DB
}

func (m *UserCategoriesModel) InsertUserCategories(userID, categoryID int) error {
	// prep the query
	q := `insert into user_categories (user_id, category_id) values ($1, $2) returning user_id`

	err := m.DB.QueryRow(q, userID, categoryID).Scan(&userID)
	return err
}

func (m *UserCategoriesModel) DeleteUserCategories(userID, categoryID int) {
	// prep the query
	q := `DELETE FROM user_categories WHERE user_id = $1 AND category_id = $2`

	m.DB.QueryRow(q, userID, categoryID)

}

func (m *UserCategoriesModel) UserCategoriesGet(userId int) ([]*Category, error) {
	q := `SELECT id, name FROM categories 
	JOIN user_categories ON categories.id = user_categories.category_id
	WHERE user_categories.user_id = $1
	ORDER BY user_categories.category_id`

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, q, userId)
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
