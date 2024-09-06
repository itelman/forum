package sqlite

import (
	"database/sql"
	"forum/internal/repository/models"
)

type ImageModel struct {
	DB *sql.DB
}

func NewImageModel(db *sql.DB) *ImageModel {
	return &ImageModel{db}
}

func (m *ImageModel) Insert(post_id int, path string) error {
	stmt := `INSERT INTO images (post_id, path) VALUES(?, ?)`

	_, err := m.DB.Exec(stmt, post_id, path)
	if err != nil {
		return err
	}

	return nil
}

func (m *ImageModel) Get(post_id int) (*models.Image, error) {
	stmt := `SELECT id, post_id, path, uploaded FROM images WHERE post_id = ?`
	row := m.DB.QueryRow(stmt, post_id)

	s := &models.Image{}
	err := row.Scan(&s.ID, &s.PostID, &s.Path, &s.Uploaded)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}
