package post_reactions

import (
	"database/sql"
	"errors"

	"github.com/itelman/forum/internal/service/post_reactions/adapters"
	"github.com/itelman/forum/internal/service/post_reactions/domain"
)

type Service interface {
	CreatePostReaction(input *CreatePostReactionInput) error
}

type service struct {
	postReactions domain.PostReactionsRepository
	posts         domain.PostsRepository
	db            *sql.DB
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
		s.postReactions = adapters.NewPostReactionsRepositorySqlite(db)
		s.posts = adapters.NewPostsRepositorySqlite(db)
		s.db = db
	}
}

func (s *service) CreatePostReaction(input *CreatePostReactionInput) error {
	makeInsertion := true

	if _, err := s.posts.Get(domain.GetPostInput{ID: input.PostID}); errors.Is(err, domain.ErrPostNotFound) {
		return domain.ErrPostReactionsBadRequest
	} else if err != nil {
		return err
	}

	reaction, err := s.postReactions.Get(domain.GetPostReactionInput{
		PostID: input.PostID,
		UserID: input.UserID,
	})
	if err != nil && !errors.Is(err, domain.ErrPostReactionNotFound) {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	if reaction != nil {
		if err := s.postReactions.Delete(tx, domain.DeletePostReactionInput{
			ID: reaction.ID,
		}); err != nil {
			tx.Rollback()
			return err
		}

		if reaction.IsLike == input.IsLike {
			makeInsertion = false
		}
	}

	if makeInsertion {
		if err := s.postReactions.Insert(tx, domain.CreatePostReactionInput{
			PostID: input.PostID,
			UserID: input.UserID,
			IsLike: input.IsLike,
		}); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := s.posts.UpdateReactionsCount(tx, domain.UpdatePostReactionsCountInput{PostID: input.PostID}); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
