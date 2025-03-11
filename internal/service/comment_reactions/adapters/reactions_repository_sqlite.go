package adapters

import (
	"database/sql"
	"errors"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/comment_reactions/domain"
)

type CommentReactionsRepositorySqlite struct {
	db *sql.DB
}

func NewCommentReactionsRepositorySqlite(db *sql.DB) *CommentReactionsRepositorySqlite {
	return &CommentReactionsRepositorySqlite{db}
}

func (r *CommentReactionsRepositorySqlite) Get(input domain.GetCommentReactionInput) (*dto.CommentReaction, error) {
	query := "SELECT id, comment_id, user_id, is_like, created FROM comment_reactions WHERE comment_id = ? AND user_id = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	reaction := &dto.CommentReaction{User: &dto.User{}}
	if err := stmt.QueryRow(input.CommentID, input.UserID).Scan(
		&reaction.ID,
		&reaction.CommentID,
		&reaction.User.ID,
		&reaction.IsLike,
		&reaction.Created,
	); errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrCommentReactionNotFound
	} else if err != nil {
		return nil, err
	}

	return reaction, nil
}

func (r *CommentReactionsRepositorySqlite) Insert(tx *sql.Tx, input domain.CreateCommentReactionInput) error {
	query := "INSERT INTO comment_reactions (comment_id, user_id, is_like) VALUES(?, ?, ?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(input.CommentID, input.UserID, input.IsLike)
	if err != nil {
		return err
	}

	return err
}

func (r *CommentReactionsRepositorySqlite) Delete(tx *sql.Tx, input domain.DeleteCommentReactionInput) error {
	query := "DELETE FROM comment_reactions WHERE id = ?"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(input.ID)
	if err != nil {
		return err
	}

	return nil
}
