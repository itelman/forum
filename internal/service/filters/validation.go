package filters

import (
	"github.com/itelman/forum/internal/service/filters/domain"
	"github.com/itelman/forum/pkg/validator"
)

type GetPostsByFiltersInput struct {
	CategoryID int
	Created    bool
	Liked      bool
	AuthUserID int
	Errors     validator.Errors
}

func (i *GetPostsByFiltersInput) validate() error {
	if i.CategoryID == -1 && !i.Created && !i.Liked {
		return domain.ErrFiltersNoneSelected

		//i.Errors.Add("generic", "Please select at least one filter")
		//return nil
	}

	if (i.Created || i.Liked) && i.AuthUserID == -1 {
		return domain.ErrUserUnauthorized
	}

	return nil
}
