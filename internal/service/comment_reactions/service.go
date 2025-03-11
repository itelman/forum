package comment_reactions

import (
	"database/sql"
	"errors"
	"github.com/itelman/forum/internal/service/comment_reactions/adapters"
	"github.com/itelman/forum/internal/service/comment_reactions/domain"
)

type Service interface {
	CreateCommentReaction(input *CreateCommentReactionInput) error
}

type service struct {
	commentReactions domain.CommentReactionsRepository
	comments         domain.CommentsRepository
	db               *sql.DB
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
		s.commentReactions = adapters.NewCommentReactionsRepositorySqlite(db)
		s.comments = adapters.NewCommentsRepositorySqlite(db)
		s.db = db
	}
}

func (s *service) CreateCommentReaction(input *CreateCommentReactionInput) error {
	makeInsertion := true

	reaction, err := s.commentReactions.Get(domain.GetCommentReactionInput{
		CommentID: input.CommentID,
		UserID:    input.UserID,
	})
	if err != nil && !errors.Is(err, domain.ErrCommentReactionNotFound) {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	if reaction != nil {
		if err := s.commentReactions.Delete(tx, domain.DeleteCommentReactionInput{
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
		if err := s.commentReactions.Insert(tx, domain.CreateCommentReactionInput{
			CommentID: input.CommentID,
			UserID:    input.UserID,
			IsLike:    input.IsLike,
		}); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := s.comments.UpdateReactionsCount(tx, domain.UpdateCommentReactionsCountInput{CommentID: input.CommentID}); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
