package posts

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/posts/adapters"
	"github.com/itelman/forum/internal/service/posts/domain"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

type Service interface {
	CreatePost(input *CreatePostInput, dir string) (*CreatePostResponse, error)
	GetPost(input *GetPostInput) (*GetPostResponse, error)
	GetAllLatestPosts() (*GetAllPostsResponse, error)
	UpdatePost(input *UpdatePostInput) error
	DeletePost(input *DeletePostInput, dir string) error
}

type service struct {
	posts          domain.PostsRepository
	postCategories domain.PostCategoriesRepository
	comments       domain.CommentsRepository
	images         domain.ImagesRepository
	db             *sql.DB
}

func NewService(opts ...Option) *service {
	svc := &service{}
	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

type Option func(*service)

func WithSqlite(db *sql.DB) Option {
	return func(s *service) {
		s.posts = adapters.NewPostsRepositorySqlite(db)
		s.postCategories = adapters.NewPostCategoriesRepositorySqlite(db)
		s.comments = adapters.NewCommentsRepositorySqlite(db)
		s.images = adapters.NewImagesRepositorySqlite(db)
		s.db = db
	}
}

type CreatePostResponse struct {
	PostID int
}

func (s *service) CreatePost(input *CreatePostInput, dir string) (*CreatePostResponse, error) {
	fileExists := false
	if !(input.ImageFile == nil || input.ImageHeader == nil) {
		fileExists = true
	}

	catgsId, err := input.validate(fileExists)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	postId, err := s.posts.Create(tx, domain.CreatePostInput{
		UserID:  input.UserID,
		Title:   input.Title,
		Content: input.Content,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := s.postCategories.Create(tx, domain.CreatePostCategoriesInput{
		PostID:       postId,
		CategoriesID: catgsId,
	}); errors.Is(err, domain.ErrPostsBadRequest) {
		input.Errors.Add("categories", "Please provide valid categories")
		tx.Rollback()
		return nil, err
	} else if err != nil {
		tx.Rollback()
		return nil, err
	}

	if fileExists {
		dirPath := filepath.Join(dir, strconv.Itoa(postId))
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			tx.Rollback()
			return nil, err
		}

		file, err := os.OpenFile(filepath.Join(dirPath, input.ImageHeader.Filename), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if _, err := io.Copy(file, input.ImageFile); err != nil {
			tx.Rollback()
			return nil, err
		}

		file.Close()

		if err := s.images.Create(tx, domain.CreateImageInput{
			PostID: postId,
			Path:   fmt.Sprintf("/images/%d/%s", postId, input.ImageHeader.Filename),
		}); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &CreatePostResponse{PostID: postId}, nil
}

type GetPostResponse struct {
	Post *dto.Post
}

func (s *service) GetPost(input *GetPostInput) (*GetPostResponse, error) {
	post, err := s.posts.Get(domain.GetPostInput{
		ID:         input.ID,
		AuthUserID: input.AuthUserID,
	})
	if err != nil {
		return nil, err
	}

	categories, err := s.postCategories.GetAllForPost(domain.GetPostCategoriesInput{PostID: input.ID})
	if err != nil {
		return nil, err
	}
	post.Categories = categories

	comments, err := s.comments.GetAllForPost(domain.GetAllCommentsForPostInput{
		PostID:         input.ID,
		AuthUserID:     input.AuthUserID,
		SortedByNewest: true,
	})
	if err != nil {
		return nil, err
	}
	post.Comments = comments

	image, err := s.images.Get(domain.GetImageInput{PostID: input.ID})
	if err != nil && !errors.Is(err, domain.ErrImageNotFound) {
		return nil, err
	}
	post.Image = image

	return &GetPostResponse{post}, nil
}

type GetAllPostsResponse struct {
	Posts []*dto.Post
}

func (s *service) GetAllLatestPosts() (*GetAllPostsResponse, error) {
	posts, err := s.posts.GetAll(domain.GetAllPostsInput{
		SortedByNewest: true,
	})
	if err != nil {
		return nil, err
	}

	return &GetAllPostsResponse{posts}, nil
}

func (s *service) UpdatePost(input *UpdatePostInput) error {
	if err := input.validate(); err != nil {
		return err
	}

	if err := s.posts.Update(domain.UpdatePostInput{
		ID:      input.ID,
		Title:   input.Title,
		Content: input.Content,
	}); err != nil {
		return err
	}

	return nil
}

func (s *service) DeletePost(input *DeletePostInput, dir string) error {
	if _, err := s.GetPost(&GetPostInput{ID: input.ID}); err != nil {
		return err
	}

	if err := os.RemoveAll(filepath.Join(dir, strconv.Itoa(input.ID))); err != nil {
		return err
	}

	if err := s.posts.Delete(domain.DeletePostInput{
		ID: input.ID,
	}); err != nil {
		return err
	}

	return nil
}
