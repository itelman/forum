package comment_reactions

import (
	"database/sql"
	"errors"

	"github.com/itelman/forum/internal/service/comment_reactions/adapters"
	"github.com/itelman/forum/internal/service/comment_reactions/domain"
)

type Service interface {
	CreateCommentReaction(input *CreateCommentReactionInput) (*CreateCommentReactionResponse, error)
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

type CreateCommentReactionResponse struct {
	PostID int
}

func (s *service) CreateCommentReaction(input *CreateCommentReactionInput) (*CreateCommentReactionResponse, error) {
	makeInsertion := true

	comment, err := s.comments.Get(domain.GetCommentInput{ID: input.CommentID})
	if errors.Is(err, domain.ErrCommentNotFound) {
		return nil, domain.ErrCommentReactionsBadRequest
	} else if err != nil {
		return nil, err
	}

	reaction, err := s.commentReactions.Get(domain.GetCommentReactionInput{
		CommentID: input.CommentID,
		UserID:    input.UserID,
	})
	if err != nil && !errors.Is(err, domain.ErrCommentReactionNotFound) {
		return nil, err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	if reaction != nil {
		if err := s.commentReactions.Delete(tx, domain.DeleteCommentReactionInput{
			ID: reaction.ID,
		}); err != nil {
			tx.Rollback()
			return nil, err
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
			return nil, err
		}
	}

	if err := s.comments.UpdateReactionsCount(tx, domain.UpdateCommentReactionsCountInput{CommentID: input.CommentID}); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &CreateCommentReactionResponse{PostID: comment.PostID}, nil
}
