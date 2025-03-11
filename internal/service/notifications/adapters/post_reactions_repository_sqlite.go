package adapters

import (
	"database/sql"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/notifications/domain"
)

type PostReactionsRepositorySqlite struct {
	db *sql.DB
}

func NewPostReactionsRepositorySqlite(db *sql.DB) *PostReactionsRepositorySqlite {
	return &PostReactionsRepositorySqlite{db}
}

func (r *PostReactionsRepositorySqlite) GetAllNotifications(input domain.GetAllPostReactionNotificationsInput) ([]*dto.PostReaction, error) {
	query := "SELECT r.post_id, u.id, u.username, r.is_like, r.created FROM post_reactions r INNER JOIN users u ON r.user_id = u.id JOIN posts p ON r.post_id = p.id WHERE p.user_id = ? AND r.user_id != ?"
	if input.SortedByNewest {
		query += " ORDER BY r.created DESC"
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

	var reactions []*dto.PostReaction
	for rows.Next() {
		reaction := &dto.PostReaction{User: &dto.User{}}

		if err := rows.Scan(
			&reaction.PostID,
			&reaction.User.ID,
			&reaction.User.Username,
			&reaction.IsLike,
			&reaction.Created,
		); err != nil {
			return nil, err
		}

		reactions = append(reactions, reaction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reactions, nil
}
