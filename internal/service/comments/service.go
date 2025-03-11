package comments

import (
	"database/sql"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/comments/adapters"
	"github.com/itelman/forum/internal/service/comments/domain"
)

type Service interface {
	CreateComment(input *CreateCommentInput) error
	GetComment(input *GetCommentInput) (*GetCommentResponse, error)
	UpdateComment(input *UpdateCommentInput) error
	DeleteComment(input *DeleteCommentInput) error
}

type service struct {
	comments domain.CommentsRepository
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
		s.comments = adapters.NewCommentsRepositorySqlite(db)
	}
}

func (s *service) CreateComment(input *CreateCommentInput) error {
	if err := input.validate(); err != nil {
		return err
	}

	if err := s.comments.Create(domain.CreateCommentInput{
		PostID:  input.PostID,
		UserID:  input.UserID,
		Content: input.Content,
	}); err != nil {
		return err
	}

	return nil
}

type GetCommentResponse struct {
	Comment *dto.Comment
}

func (s *service) GetComment(input *GetCommentInput) (*GetCommentResponse, error) {
	comment, err := s.comments.Get(domain.GetCommentInput{
		ID: input.ID,
	})
	if err != nil {
		return nil, err
	}

	return &GetCommentResponse{comment}, nil
}

type GetAllCommentsForPostResponse struct {
	Comments []*dto.Comment
}

func (s *service) UpdateComment(input *UpdateCommentInput) error {
	if err := input.validate(); err != nil {
		return err
	}

	if err := s.comments.Update(domain.UpdateCommentInput{
		ID:      input.ID,
		Content: input.Content,
	}); err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteComment(input *DeleteCommentInput) error {
	if err := s.comments.Delete(domain.DeleteCommentInput{
		ID: input.ID,
	}); err != nil {
		return err
	}

	return nil
}
