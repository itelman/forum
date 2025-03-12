package activity

import (
	"database/sql"

	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/activity/adapters"
	"github.com/itelman/forum/internal/service/activity/domain"
)

type Service interface {
	GetAllCreatedPosts(input *GetAllCreatedPostsInput) (*GetAllCreatedPostsResponse, error)
	GetAllReactedPosts(input *GetAllReactedPostsInput) (*GetAllReactedPostsResponse, error)
	GetAllCommentedPosts(input *GetAllCommentedPostsInput) (*GetAllCommentedPostsResponse, error)
}

type service struct {
	posts    domain.PostsRepository
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
		s.posts = adapters.NewPostsRepositorySqlite(db)
		s.comments = adapters.NewCommentsRepositorySqlite(db)
	}
}

type GetAllCreatedPostsResponse struct {
	Posts []*dto.Post
}

func (s *service) GetAllCreatedPosts(input *GetAllCreatedPostsInput) (*GetAllCreatedPostsResponse, error) {
	posts, err := s.posts.GetAllCreated(domain.GetAllCreatedPostsInput{
		AuthUserID:     input.AuthUserID,
		SortedByNewest: true,
	})
	if err != nil {
		return nil, err
	}

	return &GetAllCreatedPostsResponse{Posts: posts}, nil
}

type GetAllReactedPostsResponse struct {
	Posts []*dto.Post
}

func (s *service) GetAllReactedPosts(input *GetAllReactedPostsInput) (*GetAllReactedPostsResponse, error) {
	posts, err := s.posts.GetAllReacted(domain.GetAllReactedPostsInput{
		AuthUserID:     input.AuthUserID,
		SortedByNewest: true,
	})
	if err != nil {
		return nil, err
	}

	return &GetAllReactedPostsResponse{Posts: posts}, nil
}

type GetAllCommentedPostsResponse struct {
	Posts []*dto.Post
}

func (s *service) GetAllCommentedPosts(input *GetAllCommentedPostsInput) (*GetAllCommentedPostsResponse, error) {
	posts, err := s.posts.GetAllCommented(domain.GetAllCommentedPostsInput{
		AuthUserID:     input.AuthUserID,
		SortedByNewest: true,
	})
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		comments, err := s.comments.GetAllForPostByUser(domain.GetAllCommentsForPostByUserInput{
			PostID:         post.ID,
			AuthUserID:     input.AuthUserID,
			SortedByNewest: true,
		})
		if err != nil {
			return nil, err
		}

		post.Comments = comments
	}

	return &GetAllCommentedPostsResponse{Posts: posts}, nil
}
