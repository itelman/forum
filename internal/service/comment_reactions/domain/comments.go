package domain

import (
	"database/sql"
	"errors"

	"github.com/itelman/forum/internal/dto"
)

type CommentsRepository interface {
	Get(input GetCommentInput) (*dto.Comment, error)
	UpdateReactionsCount(tx *sql.Tx, input UpdateCommentReactionsCountInput) error
}

type UpdateCommentReactionsCountInput struct {
	CommentID int
}

type GetCommentInput struct {
	ID int
}

var (
	ErrCommentsBadRequest = errors.New("COMMENTS: bad request")
	ErrCommentNotFound    = errors.New("DATABASE: Comment not found")
)
