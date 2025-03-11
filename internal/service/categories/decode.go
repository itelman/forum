package categories

import (
	"github.com/itelman/forum/internal/service/categories/domain"
	"github.com/itelman/forum/pkg/validator"
	"net/http"
	"strconv"
)

func DecodeCreateCategory(r *http.Request) (interface{}, error) {
	if err := r.ParseForm(); err != nil {
		return nil, domain.ErrCategoriesBadRequest
	}

	return &CreateCategoryInput{
		Name:   r.PostForm.Get("name"),
		Errors: make(validator.Errors),
	}, nil
}

func DecodeDeleteCategory(r *http.Request) (interface{}, error) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return nil, domain.ErrCategoriesBadRequest
	}

	return &DeleteCategoryInput{ID: id}, nil
}
