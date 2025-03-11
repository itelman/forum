package domain

import (
	"database/sql"
	"errors"
	"github.com/itelman/forum/internal/dto"
)

type ImagesRepository interface {
	Create(tx *sql.Tx, input CreateImageInput) error
	Get(input GetImageInput) (*dto.Image, error)
}

type CreateImageInput struct {
	PostID int
	Path   string
}

type GetImageInput struct {
	PostID int
}

var (
	ErrImageNotFound = errors.New("DATABASE: Image not found")
)
