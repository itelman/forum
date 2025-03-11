package domain

import (
	"github.com/itelman/forum/internal/dto"
)

type CommentsRepository interface {
	GetAllForPostByUser(input GetAllCommentsForPostByUserInput) ([]*dto.Comment, error)
}

type GetAllCommentsForPostByUserInput struct {
	PostID         int
	AuthUserID     int
	SortedByNewest bool
}
