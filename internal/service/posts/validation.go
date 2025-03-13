package posts

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/posts/domain"
	"github.com/itelman/forum/pkg/validator"
)

const (
	titleMinLen = 5
	titleMaxLen = 50
	maxFileSize = 20 << 20 // 20MB
)

type CreatePostInput struct {
	UserID       int
	Title        string
	Content      string
	CategoriesID []string
	ImageFile    multipart.File
	ImageHeader  *multipart.FileHeader
	Errors       validator.Errors
}

func (i *CreatePostInput) validate(fileExists bool) ([]int, error) {
	if fileExists {
		i.validateImage()
	}

	i.validateTitle()
	i.validateContent()
	ids := i.validateCategoriesID()

	if len(i.Errors) != 0 || ids == nil {
		return nil, domain.ErrPostsBadRequest
	}

	return ids, nil
}

func (i *CreatePostInput) validateTitle() {
	if !(len(i.Title) >= titleMinLen && len(i.Title) <= titleMaxLen) {
		i.Errors.Add("title", validator.ErrInputLength(titleMinLen, titleMaxLen))
		return
	}

	if i.Title != strings.TrimSpace(i.Title) {
		i.Errors.Add("title", validator.ErrInputRequired("title"))
	}
}

func (i *CreatePostInput) validateContent() {
	if len(strings.TrimSpace(i.Content)) == 0 {
		i.Errors.Add("content", validator.ErrInputRequired("content"))
		return
	}

	if i.Content != strings.TrimSpace(i.Content) {
		i.Errors.Add("content", validator.ErrInputRequired("content"))
	}
}

func (i *CreatePostInput) validateCategoriesID() []int {
	if len(i.CategoriesID) == 0 {
		i.Errors.Add("categories", "Please provide valid categories")
		return nil
	}

	result := make([]int, 0)
	for _, idStr := range i.CategoriesID {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			i.Errors.Add("categories", "Please provide valid categories")
			return nil
		}
		result = append(result, id)
	}
	return result
}

func (i *CreatePostInput) validateImage() {
	if i.ImageHeader.Size > maxFileSize {
		i.Errors.Add("image", fmt.Sprintf("Max size exceeded (max %d MB)", 20))
		return
	}

	extIsValid := false

	exts := []string{".jpg", ".jpeg", ".png", ".gif"}
	for _, ext := range exts {
		if filepath.Ext(i.ImageHeader.Filename) == ext {
			extIsValid = true
			break
		}
	}

	if !extIsValid {
		i.Errors.Add("image", fmt.Sprintf("File extension should be one of the following: %s", exts))
	}
}

type GetPostInput struct {
	ID         int
	AuthUserID int
}

type UpdatePostInput struct {
	ID      int
	Title   string
	Content string
	Errors  validator.Errors
}

func (i *UpdatePostInput) validate(post *dto.Post) error {
	i.validateTitle()
	i.validateContent()

	if len(i.Errors) != 0 {
		return domain.ErrPostsBadRequest
	}

	if i.Title == post.Title && i.Content == post.Content {
		i.Errors.Add("generic", validator.ErrInputUnchanged)
		return domain.ErrPostsBadRequest
	}

	return nil
}

func (i *UpdatePostInput) validateTitle() {
	if i.Title != strings.TrimSpace(i.Title) {
		i.Errors.Add("title", validator.ErrInputRequired("title"))
		return
	}

	if !(len(i.Title) >= titleMinLen && len(i.Title) <= titleMaxLen) {
		i.Errors.Add("title", validator.ErrInputLength(titleMinLen, titleMaxLen))
		return
	}
}

func (i *UpdatePostInput) validateContent() {
	if len(strings.TrimSpace(i.Content)) == 0 {
		i.Errors.Add("content", validator.ErrInputRequired("content"))
		return
	}

	if i.Content != strings.TrimSpace(i.Content) {
		i.Errors.Add("content", validator.ErrInputRequired("content"))
		return
	}
}

type DeletePostInput struct {
	ID int
}
