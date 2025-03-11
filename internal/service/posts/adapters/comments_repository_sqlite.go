package adapters

import (
	"database/sql"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/posts/domain"
)

type CommentsRepositorySqlite struct {
	db *sql.DB
}

func NewCommentsRepositorySqlite(db *sql.DB) *CommentsRepositorySqlite {
	return &CommentsRepositorySqlite{db}
}

func (r *CommentsRepositorySqlite) GetAllForPost(input domain.GetAllCommentsForPostInput) ([]*dto.Comment, error) {
	query := "SELECT comments.id, comments.post_id, users.id, users.username, comments.content, comments.likes, comments.dislikes, comments.created, COALESCE(comment_reactions.is_like, -1) AS is_like FROM comments INNER JOIN users ON comments.user_id = users.id LEFT JOIN comment_reactions ON comments.id = comment_reactions.comment_id AND (comment_reactions.user_id = ? OR ? = -1) WHERE comments.post_id = ?"
	if input.SortedByNewest {
		query += " ORDER BY comments.created DESC"
	}
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(input.AuthUserID, input.AuthUserID, input.PostID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*dto.Comment{}
	for rows.Next() {
		comment := &dto.Comment{User: &dto.User{}}

		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.User.ID,
			&comment.User.Username,
			&comment.Content,
			&comment.Likes,
			&comment.Dislikes,
			&comment.Created,
			&comment.AuthUserReaction,
		); err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
