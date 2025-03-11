package filters

import (
	"database/sql"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/filters/adapters"
	"github.com/itelman/forum/internal/service/filters/domain"
)

type Service interface {
	GetPostsByFilters(input *GetPostsByFiltersInput) (*GetManyPostsResponse, error)
}

type service struct {
	posts domain.PostsRepository
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
	}
}

type GetManyPostsResponse struct {
	Posts []*dto.Post
}

func (s *service) GetPostsByFilters(input *GetPostsByFiltersInput) (*GetManyPostsResponse, error) {
	if err := input.validate(); err != nil {
		return nil, err
	}
	
	posts, err := s.posts.GetManyByFilters(domain.GetPostsByFiltersInput{
		CategoryID:     input.CategoryID,
		Created:        input.Created,
		Liked:          input.Liked,
		AuthUserID:     input.AuthUserID,
		SortedByNewest: true,
	})
	if err != nil {
		return nil, err
	}

	return &GetManyPostsResponse{Posts: posts}, nil
}
