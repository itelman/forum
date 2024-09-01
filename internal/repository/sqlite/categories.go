package sqlite

import (
	"database/sql"
	"forum/internal/repository/models"
)

type CategoryModel struct {
	DB *sql.DB
}

func NewCategoryModel(db *sql.DB) *CategoryModel {
	return &CategoryModel{db}
}

func (m *CategoryModel) Get(id int) (*models.Category, error) {
	stmt := `SELECT id, name, created FROM categories WHERE id = ?`
	row := m.DB.QueryRow(stmt, id)

	s := &models.Category{}
	err := row.Scan(&s.ID, &s.Name, &s.Created)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

func (m *CategoryModel) Latest() ([]*models.Category, error) {
	stmt := `SELECT * FROM categories`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []*models.Category{}

	for rows.Next() {
		s := &models.Category{}
		err := rows.Scan(&s.ID, &s.Name, &s.Created)
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
