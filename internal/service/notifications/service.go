package notifications

import (
	"database/sql"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/notifications/adapters"
	"github.com/itelman/forum/internal/service/notifications/domain"
)

type Service interface {
	GetAllCommentNotifications(input *GetAllCommentNotificationsInput) (*GetAllCommentNotificationsResponse, error)
	GetAllPostReactionNotifications(input *GetAllPostReactionNotificationsInput) (*GetAllPostReactionNotificationsResponse, error)
}

type service struct {
	comments      domain.CommentsRepository
	postReactions domain.PostReactionsRepository
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
		s.postReactions = adapters.NewPostReactionsRepositorySqlite(db)
	}
}

type GetAllCommentNotificationsResponse struct {
	Comments []*dto.Comment
}

func (s *service) GetAllCommentNotifications(input *GetAllCommentNotificationsInput) (*GetAllCommentNotificationsResponse, error) {
	comments, err := s.comments.GetAllNotifications(domain.GetAllCommentNotificationsInput{
		AuthUserID:     input.AuthUserID,
		SortedByNewest: true,
	})
	if err != nil {
		return nil, err
	}

	return &GetAllCommentNotificationsResponse{comments}, nil
}

type GetAllPostReactionNotificationsResponse struct {
	PostReactions []*dto.PostReaction
}

func (s *service) GetAllPostReactionNotifications(input *GetAllPostReactionNotificationsInput) (*GetAllPostReactionNotificationsResponse, error) {
	reactions, err := s.postReactions.GetAllNotifications(domain.GetAllPostReactionNotificationsInput{
		AuthUserID:     input.AuthUserID,
		SortedByNewest: true,
	})
	if err != nil {
		return nil, err
	}

	return &GetAllPostReactionNotificationsResponse{reactions}, nil
}
