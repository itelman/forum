package mock

import (
	"forum/internal/repository/models"
)

type ImageModel struct {
	DB []*models.Image
}

func NewImageModel() *ImageModel {
	return &ImageModel{make([]*models.Image, 0)}
}

func (m *ImageModel) Insert(post_id int, path string) error {
	return nil
}

func (m *ImageModel) Get(post_id int) (*models.Image, error) {
	s := &models.Image{}

	return s, nil
}
