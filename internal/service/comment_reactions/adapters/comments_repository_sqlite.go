package adapters

import (
	"database/sql"
	"github.com/itelman/forum/internal/service/comment_reactions/domain"
)

type CommentsRepositorySqlite struct {
	db *sql.DB
}

func NewCommentsRepositorySqlite(db *sql.DB) *CommentsRepositorySqlite {
	return &CommentsRepositorySqlite{db}
}

func (r *CommentsRepositorySqlite) UpdateReactionsCount(tx *sql.Tx, input domain.UpdateCommentReactionsCountInput) error {
	query := "UPDATE comments SET likes = (SELECT COUNT(*) FROM comment_reactions WHERE comment_id = comments.id AND is_like = 1), dislikes = (SELECT COUNT(*) FROM comment_reactions WHERE comment_id = comments.id AND is_like = 0) WHERE id = ?"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(input.CommentID); err != nil {
		return err
	}

	return nil
}
