package sqlite

import (
	"database/sql"
	"forum/internal/repository/models"

	"github.com/mattn/go-sqlite3"
)

type PostReactionModel struct {
	DB *sql.DB
}

func NewPostReactionModel(db *sql.DB) *PostReactionModel {
	return &PostReactionModel{db}
}

func (m *PostReactionModel) Insert(post_id, user_id, is_like int) error {
	existingLike, err := m.Get(post_id, user_id)
	if err != nil {
		return err
	}

	if existingLike == is_like {
		if err := m.Delete(post_id, user_id); err != nil {
			return err
		}
		return nil
	} else if existingLike != -1 {
		stmt := `UPDATE post_reactions SET is_like = ? WHERE post_id = ? AND user_id = ?`
		_, err = m.DB.Exec(stmt, is_like, post_id, user_id)
		if err != nil {
			return err
		}
		return nil
	}

	stmt := `INSERT INTO post_reactions (post_id, user_id, is_like) VALUES(?, ?, ?)`
	_, err = m.DB.Exec(stmt, post_id, user_id, is_like)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint && sqliteErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
			return models.ErrNoRecord
		}
		return err
	}

	return nil
}

func (m *PostReactionModel) Delete(post_id, user_id int) error {
	stmt := `DELETE FROM post_reactions WHERE post_id = $1 AND user_id = $2`
	_, err := m.DB.Exec(stmt, post_id, user_id)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostReactionModel) Get(post_id, user_id int) (int, error) {
	var isLike int

	stmt := `SELECT is_like FROM post_reactions WHERE post_id = $1 AND user_id = $2`
	err := m.DB.QueryRow(stmt, post_id, user_id).Scan(&isLike)
	if err == sql.ErrNoRows {
		return -1, nil
	} else if err != nil {
		return -1, err
	}

	return isLike, nil
}

func (m *PostReactionModel) Likes(post_id int) (int, error) {
	var count int
	err := m.DB.QueryRow("SELECT COUNT(*) FROM post_reactions WHERE post_id = $1 AND is_like = 1", post_id).Scan(&count)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return -1, err
	}
	return count, nil
}

func (m *PostReactionModel) Dislikes(post_id int) (int, error) {
	var count int
	err := m.DB.QueryRow("SELECT COUNT(*) FROM post_reactions WHERE post_id = $1 AND is_like = 0", post_id).Scan(&count)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return -1, err
	}

	return count, nil
}

func (m *PostReactionModel) LatestIgnoreUser(post_id, user_id int) ([]*models.PostReaction, error) {
	stmt := `SELECT post_reactions.id, post_reactions.post_id, users.id, users.name, post_reactions.is_like, post_reactions.created FROM post_reactions INNER JOIN users ON post_reactions.user_id = users.id WHERE post_reactions.post_id = ? AND post_reactions.user_id <> ? ORDER BY post_reactions.created DESC`
	rows, err := m.DB.Query(stmt, post_id, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	post_reactions := []*models.PostReaction{}

	for rows.Next() {
		s := models.PostReaction{}
		err := rows.Scan(&s.ID, &s.PostID, &s.UserID, &s.Username, &s.IsLike, &s.Created)
		if err != nil {
			return nil, err
		}
		post_reactions = append(post_reactions, &s)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return post_reactions, nil
}

func (m *PostReactionModel) FilterByLiked(post_id, user_id int, val string) (bool, error) {
	if len(val) == 0 || val == "0" || user_id == -1 {
		return true, nil
	}

	stmt := `SELECT COUNT(*) FROM post_reactions WHERE post_id = ? AND user_id = ? AND is_like = 1`

	var count int
	err := m.DB.QueryRow(stmt, post_id, user_id).Scan(&count)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (m *PostReactionModel) GetReactionsByUser(user_id int) ([]*models.PostReaction, error) {
	stmt := `SELECT id, post_id, user_id, is_like, created FROM post_reactions WHERE user_id = ? ORDER BY created DESC`
	rows, err := m.DB.Query(stmt, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	post_reactions := []*models.PostReaction{}

	for rows.Next() {
		s := models.PostReaction{}
		err := rows.Scan(&s.ID, &s.PostID, &s.UserID, &s.IsLike, &s.Created)
		if err != nil {
			return nil, err
		}
		post_reactions = append(post_reactions, &s)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return post_reactions, nil
}
