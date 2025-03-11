package domain

import (
	"errors"
	"github.com/itelman/forum/internal/dto"
)

type CategoriesRepository interface {
	Create(input CreateCategoryInput) error
	Get(input GetCategoryInput) (*dto.Category, error)
	GetAll(input GetAllCategoriesInput) ([]*dto.Category, error)
	Delete(input DeleteCategoryInput) error
}

type CreateCategoryInput struct {
	Name string
}

type GetCategoryInput struct {
	ID int
}

type GetAllCategoriesInput struct {
	SortedByNewest bool
}

type DeleteCategoryInput struct {
	ID int
}

var (
	ErrCategoriesBadRequest = errors.New("CATEGORIES: bad request")
	ErrCategoryNotFound     = errors.New("DATABASE: Category not found")
)
