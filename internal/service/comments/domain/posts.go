package domain

import (
	"errors"

	"github.com/itelman/forum/internal/dto"
)

type PostsRepository interface {
	Get(input GetPostInput) (*dto.Post, error)
}

type GetPostInput struct {
	ID int
}

var (
	ErrPostsBadRequest = errors.New("POSTS: bad request")
	ErrPostNotFound    = errors.New("DATABASE: Post not found")
)
