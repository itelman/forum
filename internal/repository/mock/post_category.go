package mock

import (
	"forum/internal/repository/models"
)

type PostCategoryModel struct {
	DB []*models.PostCategory
}

func NewPostCategoryModel() *PostCategoryModel {
	return &PostCategoryModel{make([]*models.PostCategory, 0)}
}

func (m *PostCategoryModel) Insert(post_id string, categories_id []string) error {
	return nil
}

func (m *PostCategoryModel) Get(id int) ([]string, error) {
	var categories []string

	return categories, nil
}

func (m *PostCategoryModel) FilterByCategories(post_id int, categories_id []string) (bool, error) {
	return false, nil
}
