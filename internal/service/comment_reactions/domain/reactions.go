package domain

import (
	"database/sql"
	"errors"
	"github.com/itelman/forum/internal/dto"
)

type CommentReactionsRepository interface {
	Get(input GetCommentReactionInput) (*dto.CommentReaction, error)
	Insert(tx *sql.Tx, input CreateCommentReactionInput) error
	Delete(tx *sql.Tx, input DeleteCommentReactionInput) error
}

type CreateCommentReactionInput struct {
	CommentID int
	UserID    int
	IsLike    int
}

type GetCommentReactionInput struct {
	CommentID int
	UserID    int
}

type DeleteCommentReactionInput struct {
	ID int
}

var (
	ErrCommentReactionsBadRequest = errors.New("COMMENT REACTIONS: bad request")
	ErrCommentReactionNotFound    = errors.New("DATABASE: Post reaction not found")
)
