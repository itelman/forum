package domain

import (
	"github.com/itelman/forum/internal/dto"
)

type CommentsRepository interface {
	GetAllForPost(input GetAllCommentsForPostInput) ([]*dto.Comment, error)
}

type GetAllCommentsForPostInput struct {
	PostID         int
	AuthUserID     int
	SortedByNewest bool
}
