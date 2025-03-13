package comments

import (
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/comments/domain"
	"github.com/itelman/forum/pkg/validator"
	"net/http"
	"strconv"
)

func DecodeCreateComment(r *http.Request) (interface{}, error) {
	if err := r.ParseForm(); err != nil {
		return nil, domain.ErrCommentsBadRequest
	}

	postId, err := strconv.Atoi(r.PostForm.Get("post_id"))
	if err != nil {
		return nil, domain.ErrCommentsBadRequest
	}

	return &CreateCommentInput{
		PostID:  postId,
		UserID:  dto.GetAuthUser(r).ID,
		Content: r.PostForm.Get("content"),
		Errors:  make(validator.Errors),
	}, nil
}

func DecodeGetComment(r *http.Request) (interface{}, error) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return nil, domain.ErrCommentsBadRequest
	}

	return &GetCommentInput{
		ID: id,
	}, nil
}

func DecodeUpdateComment(r *http.Request) (interface{}, error) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return nil, domain.ErrCommentsBadRequest
	}

	if err := r.ParseForm(); err != nil {
		return nil, domain.ErrCommentsBadRequest
	}

	return &UpdateCommentInput{
		ID:      id,
		Content: r.PostForm.Get("content"),
		Errors:  make(validator.Errors),
	}, nil
}

func DecodeDeleteComment(r *http.Request) (interface{}, error) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return nil, domain.ErrCommentsBadRequest
	}

	return &DeleteCommentInput{
		ID: id,
	}, nil
}
