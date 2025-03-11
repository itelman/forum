package adapters

import (
	"database/sql"
	"errors"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/posts/domain"
)

type ImagesRepositorySqlite struct {
	db *sql.DB
}

func NewImagesRepositorySqlite(db *sql.DB) *ImagesRepositorySqlite {
	return &ImagesRepositorySqlite{db}
}

func (r *ImagesRepositorySqlite) Create(tx *sql.Tx, input domain.CreateImageInput) error {
	query := "INSERT INTO images (post_id, path) VALUES(?, ?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(input.PostID, input.Path)
	if err != nil {
		return err
	}

	return nil
}

func (r *ImagesRepositorySqlite) Get(input domain.GetImageInput) (*dto.Image, error) {
	query := "SELECT id, path, uploaded FROM images WHERE post_id = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	image := &dto.Image{}
	if err := stmt.QueryRow(input.PostID).Scan(
		&image.ID,
		&image.Path,
		&image.Uploaded,
	); errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrImageNotFound
	} else if err != nil {
		return nil, err
	}

	return image, nil
}
