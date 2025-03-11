package domain

import (
	"errors"
	"github.com/itelman/forum/internal/dto"
)

type PostsRepository interface {
	GetManyByFilters(input GetPostsByFiltersInput) ([]*dto.Post, error)
}

type GetPostsByFiltersInput struct {
	CategoryID     int
	Created        bool
	Liked          bool
	AuthUserID     int
	SortedByNewest bool
}

var (
	ErrFiltersBadRequest   = errors.New("FILTERS: bad request")
	ErrFiltersNoneSelected = errors.New("FILTERS: none selected")
	ErrUserUnauthorized    = errors.New("FILTERS: user unauthorized")
)
