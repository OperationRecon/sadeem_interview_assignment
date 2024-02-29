package data

import "database/sql"

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
