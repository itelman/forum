package activity

import (
	"github.com/itelman/forum/internal/dto"
	"net/http"
)

func DecodeGetAllCreatedPosts(r *http.Request) interface{} {
	return &GetAllCreatedPostsInput{dto.GetAuthUser(r).ID}
}

func DecodeGetAllReactedPosts(r *http.Request) interface{} {
	return &GetAllReactedPostsInput{dto.GetAuthUser(r).ID}
}

func DecodeGetAllCommentedPosts(r *http.Request) interface{} {
	return &GetAllCommentedPostsInput{dto.GetAuthUser(r).ID}
}
