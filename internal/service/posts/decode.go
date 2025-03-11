package posts

import (
	"errors"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/posts/domain"
	"github.com/itelman/forum/pkg/validator"
	"net/http"
	"strconv"
)

func DecodeCreatePost(r *http.Request) (interface{}, error) {
	if err := r.ParseMultipartForm(maxFileSize + (1 << 20)); err != nil {
		return nil, domain.ErrPostsBadRequest
	}

	file, header, err := r.FormFile("image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		return nil, err
	}

	return &CreatePostInput{
		UserID:       dto.GetAuthUser(r).ID,
		Title:        r.PostForm.Get("title"),
		Content:      r.PostForm.Get("content"),
		CategoriesID: r.PostForm["categories_id"],
		ImageFile:    file,
		ImageHeader:  header,
		Errors:       make(validator.Errors),
	}, nil
}

func DecodeGetPost(r *http.Request) (interface{}, error) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return nil, domain.ErrPostsBadRequest
	}

	userId := -1
	user := dto.GetAuthUser(r)
	if user != nil {
		userId = user.ID
	}

	return &GetPostInput{
		ID:         id,
		AuthUserID: userId,
	}, nil
}

func DecodeUpdatePost(r *http.Request) (interface{}, error) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return nil, domain.ErrPostsBadRequest
	}

	if err := r.ParseForm(); err != nil {
		return nil, domain.ErrPostsBadRequest
	}

	return &UpdatePostInput{
		ID:      id,
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Errors:  make(validator.Errors),
	}, nil
}

func DecodeDeletePost(r *http.Request) (interface{}, error) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return nil, domain.ErrPostsBadRequest
	}

	return &DeletePostInput{
		ID: id,
	}, nil
}
