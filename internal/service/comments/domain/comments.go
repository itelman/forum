package domain

import (
	"errors"
	"github.com/itelman/forum/internal/dto"
)

type CommentsRepository interface {
	Create(input CreateCommentInput) error
	Get(input GetCommentInput) (*dto.Comment, error)
	Update(input UpdateCommentInput) error
	Delete(input DeleteCommentInput) error
}

type CreateCommentInput struct {
	PostID  int
	UserID  int
	Content string
}

type GetCommentInput struct {
	ID int
}

type GetAllCommentsForPostInput struct {
	PostID         int
	AuthUserID     int
	SortedByNewest bool
}

type UpdateCommentInput struct {
	ID      int
	Content string
}

type DeleteCommentInput struct {
	ID int
}

var (
	ErrCommentsBadRequest = errors.New("COMMENTS: bad request")
	ErrCommentNotFound    = errors.New("DATABASE: Comment not found")
)
