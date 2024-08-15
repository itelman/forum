package sqlite

import (
	"database/sql"
	"fmt"
	"strings"
)

type PostCategoryModel struct {
	DB *sql.DB
}

func (m *PostCategoryModel) Insert(post_id string, categories_id []string) error {
	var stmt string

	for _, id := range categories_id {
		stmt = `INSERT INTO post_category (post_id, category_id) VALUES(?, ?)`
		_, err := m.DB.Exec(stmt, post_id, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *PostCategoryModel) Get(id int) ([]string, error) {
	stmt := `SELECT categories.name FROM post_category INNER JOIN categories ON post_category.category_id = categories.id WHERE post_category.post_id = ?`
	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string

	for rows.Next() {
		var s string
		err := rows.Scan(&s)
		if err != nil {
			return nil, err
		}
		categories = append(categories, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (m *PostCategoryModel) FilterByCategories(post_id int, categories_id []string, val int) (bool, error) {
	if val == 0 {
		return true, nil
	}

	placeholders := strings.Repeat("?,", len(categories_id))
	placeholders = placeholders[:len(placeholders)-1] 

	stmt := fmt.Sprintf("SELECT COUNT(*) FROM post_category WHERE category_id IN (%s) AND post_id = ?", placeholders)

	args := make([]interface{}, len(categories_id)+1)
	for i, id := range categories_id {
		args[i] = id
	}

	var count int
	args[len(categories_id)] = post_id

	err := m.DB.QueryRow(stmt, args...).Scan(&count)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}
