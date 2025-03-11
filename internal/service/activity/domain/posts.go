package domain

import (
	"github.com/itelman/forum/internal/dto"
)

type PostsRepository interface {
	GetAllCreated(input GetAllCreatedPostsInput) ([]*dto.Post, error)
	GetAllReacted(input GetAllReactedPostsInput) ([]*dto.Post, error)
	GetAllCommented(input GetAllCommentedPostsInput) ([]*dto.Post, error)
}

type GetAllCreatedPostsInput struct {
	AuthUserID     int
	SortedByNewest bool
}

type GetAllReactedPostsInput struct {
	AuthUserID     int
	SortedByNewest bool
}

type GetAllCommentedPostsInput struct {
	AuthUserID     int
	SortedByNewest bool
}
