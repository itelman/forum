package domain

import (
	"database/sql"
	"errors"
	"github.com/itelman/forum/internal/dto"
)

type PostsRepository interface {
	Create(tx *sql.Tx, input CreatePostInput) (int, error)
	Get(input GetPostInput) (*dto.Post, error)
	GetAll(input GetAllPostsInput) ([]*dto.Post, error)
	Update(input UpdatePostInput) error
	Delete(input DeletePostInput) error
}

type CreatePostInput struct {
	UserID  int
	Title   string
	Content string
}

type GetPostInput struct {
	ID         int
	AuthUserID int
}

type GetAllPostsInput struct {
	SortedByNewest bool
}

type UpdatePostInput struct {
	ID      int
	Title   string
	Content string
}

type DeletePostInput struct {
	ID int
}

var (
	ErrPostsBadRequest = errors.New("POSTS: bad request")
	ErrPostNotFound    = errors.New("DATABASE: Post not found")
)
