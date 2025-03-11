package categories

import (
	"github.com/itelman/forum/internal/service/categories/domain"
	"github.com/itelman/forum/pkg/validator"
	"strings"
)

type CreateCategoryInput struct {
	Name   string
	Errors validator.Errors
}

func (i *CreateCategoryInput) validate() error {
	i.validateName()

	if len(i.Errors) != 0 {
		return domain.ErrCategoriesBadRequest
	}

	i.Name = strings.ToTitle(strings.ToLower(i.Name))

	return nil
}

func (i *CreateCategoryInput) validateName() {
	if len(strings.TrimSpace(i.Name)) == 0 {
		i.Errors.Add("name", validator.ErrInputRequired("name"))
		return
	}

	if i.Name != strings.TrimSpace(i.Name) {
		i.Errors.Add("name", validator.ErrInputRequired("name"))
	}
}

type DeleteCategoryInput struct {
	ID int
}
