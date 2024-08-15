package sqlite

import (
	"database/sql"
	"forum/internal/repository/models"
)

type CategoryModel struct {
	DB *sql.DB
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
