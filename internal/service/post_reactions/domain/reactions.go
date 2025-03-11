package domain

import (
	"database/sql"
	"errors"
	"github.com/itelman/forum/internal/dto"
)

type PostReactionsRepository interface {
	Get(input GetPostReactionInput) (*dto.PostReaction, error)
	Insert(tx *sql.Tx, input CreatePostReactionInput) error
	Delete(tx *sql.Tx, input DeletePostReactionInput) error
}

type CreatePostReactionInput struct {
	PostID int
	UserID int
	IsLike int
}

type GetPostReactionInput struct {
	PostID int
	UserID int
}

type DeletePostReactionInput struct {
	ID int
}

var (
	ErrPostReactionsBadRequest = errors.New("POST REACTIONS: bad request")
	ErrPostReactionNotFound    = errors.New("DATABASE: Post reaction not found")
)
