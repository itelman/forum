package domain

import "database/sql"

type PostCategoriesRepository interface {
	Create(tx *sql.Tx, input CreatePostCategoriesInput) error
	GetAllForPost(input GetPostCategoriesInput) ([]string, error)
}

type CreatePostCategoriesInput struct {
	PostID       int
	CategoriesID []int
}

type GetPostCategoriesInput struct {
	PostID int
}
