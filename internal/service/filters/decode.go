package filters

import (
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/filters/domain"
	"github.com/itelman/forum/pkg/validator"
	"net/http"
	"strconv"
)

func DecodeGetPostsByFilters(r *http.Request) (interface{}, error) {
	if err := r.ParseForm(); err != nil {
		return nil, domain.ErrFiltersBadRequest
	}

	catgId, err := strconv.Atoi(r.PostForm.Get("category_id"))
	if err != nil {
		catgId = -1
	}

	userId := -1
	user := dto.GetAuthUser(r)
	if user != nil {
		userId = user.ID
	}

	return &GetPostsByFiltersInput{
		CategoryID: catgId,
		Created:    r.PostForm.Get("created") == "1",
		Liked:      r.PostForm.Get("liked") == "1",
		AuthUserID: userId,
		Errors:     make(validator.Errors),
	}, nil
}
