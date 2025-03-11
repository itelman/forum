package adapters

import (
	"database/sql"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/notifications/domain"
)

type CommentsRepositorySqlite struct {
	db *sql.DB
}

func NewCommentsRepositorySqlite(db *sql.DB) *CommentsRepositorySqlite {
	return &CommentsRepositorySqlite{db}
}

func (r *CommentsRepositorySqlite) GetAllNotifications(input domain.GetAllCommentNotificationsInput) ([]*dto.Comment, error) {
	query := "SELECT c.post_id, u.id, u.username, c.created FROM comments c INNER JOIN users u ON c.user_id = u.id JOIN posts p ON c.post_id = p.id WHERE p.user_id = ? AND c.user_id != ?"
	if input.SortedByNewest {
		query += " ORDER BY c.created DESC"
	}

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(input.AuthUserID, input.AuthUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*dto.Comment
	for rows.Next() {
		comment := &dto.Comment{User: &dto.User{}}

		if err := rows.Scan(
			&comment.PostID,
			&comment.User.ID,
			&comment.User.Username,
			&comment.Created,
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
