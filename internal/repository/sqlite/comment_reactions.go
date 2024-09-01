package sqlite

import (
	"database/sql"
	"fmt"
	"forum/internal/repository/models"

	"github.com/mattn/go-sqlite3"
)

type CommentReactionModel struct {
	DB *sql.DB
}

func NewCommentReactionModel(db *sql.DB) *CommentReactionModel {
	return &CommentReactionModel{db}
}

func (m *CommentReactionModel) Insert(comment_id, user_id, is_like int) error {
	existingLike, err := m.Get(comment_id, user_id)
	if err != nil {
		return err
	}

	if existingLike == is_like {
		if err := m.Delete(comment_id, user_id); err != nil {
			return err
		}
		return nil
	} else if existingLike != -1 {
		stmt := `UPDATE comment_reactions SET is_like = ? WHERE comment_id = ? AND user_id = ?`
		_, err = m.DB.Exec(stmt, is_like, comment_id, user_id)
		if err != nil {
			return err
		}
		return nil
	}
	fmt.Println(is_like)
	stmt := `INSERT INTO comment_reactions (comment_id, user_id, is_like) VALUES(?, ?, ?)`
	_, err = m.DB.Exec(stmt, comment_id, user_id, is_like)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint && sqliteErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
			return models.ErrNoRecord
		}
		return err
	}

	return nil
}

func (m *CommentReactionModel) Delete(comment_id, user_id int) error {
	stmt := `DELETE FROM comment_reactions WHERE comment_id = $1 AND user_id = $2`
	_, err := m.DB.Exec(stmt, comment_id, user_id)
	if err != nil {
		return err
	}

	return nil
}

func (m *CommentReactionModel) Get(comment_id, user_id int) (int, error) {
	var isLike int

	stmt := `SELECT is_like FROM comment_reactions WHERE comment_id = $1 AND user_id = $2`
	err := m.DB.QueryRow(stmt, comment_id, user_id).Scan(&isLike)
	if err == sql.ErrNoRows {
		return -1, nil
	} else if err != nil {
		return -1, err
	}
	return isLike, nil
}

func (m *CommentReactionModel) Likes(comment_id int) (int, error) {
	var count int
	err := m.DB.QueryRow("SELECT COUNT(*) FROM comment_reactions WHERE comment_id = $1 AND is_like = 1", comment_id).Scan(&count)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return -1, err
	}

	return count, nil
}

func (m *CommentReactionModel) Dislikes(comment_id int) (int, error) {
	var count int
	err := m.DB.QueryRow("SELECT COUNT(*) FROM comment_reactions WHERE comment_id = $1 AND is_like = 0", comment_id).Scan(&count)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return -1, err
	}

	return count, nil
}
