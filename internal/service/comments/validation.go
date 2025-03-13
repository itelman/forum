package comments

import (
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/comments/domain"
	"github.com/itelman/forum/pkg/validator"
	"strings"
)

type CreateCommentInput struct {
	PostID  int
	UserID  int
	Content string
	Errors  validator.Errors
}

func (i *CreateCommentInput) validate() error {
	i.validateContent()

	if len(i.Errors) != 0 {
		return domain.ErrCommentsBadRequest
	}

	return nil
}

func (i *CreateCommentInput) validateContent() {
	if len(strings.TrimSpace(i.Content)) == 0 {
		i.Errors.Add("content", validator.ErrInputRequired("content"))
		return
	}

	if i.Content != strings.TrimSpace(i.Content) {
		i.Errors.Add("content", validator.ErrInputRequired("content"))
	}
}

type GetCommentInput struct {
	ID int
}

type UpdateCommentInput struct {
	ID      int
	Content string
	Errors  validator.Errors
}

func (i *UpdateCommentInput) validate(comment *dto.Comment) error {
	i.validateContent(comment)

	if len(i.Errors) != 0 {
		return domain.ErrCommentsBadRequest
	}

	return nil
}

func (i *UpdateCommentInput) validateContent(comment *dto.Comment) {
	if len(strings.TrimSpace(i.Content)) == 0 {
		i.Errors.Add("content", validator.ErrInputRequired("content"))
		return
	}

	if i.Content != strings.TrimSpace(i.Content) {
		i.Errors.Add("content", validator.ErrInputRequired("content"))
		return
	}

	if i.Content == comment.Content {
		i.Errors.Add("content", validator.ErrInputUnchanged)
	}
}

type DeleteCommentInput struct {
	ID int
}
