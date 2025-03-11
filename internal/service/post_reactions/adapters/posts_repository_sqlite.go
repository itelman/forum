package adapters

import (
	"database/sql"
	"github.com/itelman/forum/internal/service/post_reactions/domain"
)

type PostsRepositorySqlite struct {
	db *sql.DB
}

func NewPostsRepositorySqlite(db *sql.DB) *PostsRepositorySqlite {
	return &PostsRepositorySqlite{db}
}

func (r *PostsRepositorySqlite) UpdateReactionsCount(tx *sql.Tx, input domain.UpdatePostReactionsCountInput) error {
	query := "UPDATE posts SET likes = (SELECT COUNT(*) FROM post_reactions WHERE post_id = posts.id AND is_like = 1), dislikes = (SELECT COUNT(*) FROM post_reactions WHERE post_id = posts.id AND is_like = 0) WHERE id = ?"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(input.PostID); err != nil {
		return err
	}

	return nil
}
