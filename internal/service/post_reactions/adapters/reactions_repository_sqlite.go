package adapters

import (
	"database/sql"
	"errors"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/post_reactions/domain"
)

type PostReactionsRepositorySqlite struct {
	db *sql.DB
}

func NewPostReactionsRepositorySqlite(db *sql.DB) *PostReactionsRepositorySqlite {
	return &PostReactionsRepositorySqlite{db}
}

func (r *PostReactionsRepositorySqlite) Get(input domain.GetPostReactionInput) (*dto.PostReaction, error) {
	query := "SELECT id, post_id, user_id, is_like, created FROM post_reactions WHERE post_id = ? AND user_id = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	reaction := &dto.PostReaction{User: &dto.User{}}
	if err := stmt.QueryRow(input.PostID, input.UserID).Scan(
		&reaction.ID,
		&reaction.PostID,
		&reaction.User.ID,
		&reaction.IsLike,
		&reaction.Created,
	); errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrPostReactionNotFound
	} else if err != nil {
		return nil, err
	}

	return reaction, nil
}

func (r *PostReactionsRepositorySqlite) Insert(tx *sql.Tx, input domain.CreatePostReactionInput) error {
	query := "INSERT INTO post_reactions (post_id, user_id, is_like) VALUES(?, ?, ?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(input.PostID, input.UserID, input.IsLike)
	if err != nil {
		return err
	}

	return err
}

func (r *PostReactionsRepositorySqlite) Delete(tx *sql.Tx, input domain.DeletePostReactionInput) error {
	query := "DELETE FROM post_reactions WHERE id = ?"
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
