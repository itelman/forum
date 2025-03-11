package adapters

import (
	"database/sql"
	"errors"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/comments/domain"
)

type CommentsRepositorySqlite struct {
	db *sql.DB
}

func NewCommentsRepositorySqlite(db *sql.DB) *CommentsRepositorySqlite {
	return &CommentsRepositorySqlite{db}
}

func (r *CommentsRepositorySqlite) Create(input domain.CreateCommentInput) error {
	query := "INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(input.PostID, input.UserID, input.Content)
	if err != nil {
		return err
	}

	return err
}

func (r *CommentsRepositorySqlite) Get(input domain.GetCommentInput) (*dto.Comment, error) {
	query := "SELECT comments.id, comments.post_id, users.id, users.username, comments.content, comments.likes, comments.dislikes, comments.created FROM comments INNER JOIN users ON comments.user_id = users.id WHERE comments.id = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	comment := &dto.Comment{User: &dto.User{}}
	if err := stmt.QueryRow(input.ID).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.User.ID,
		&comment.User.Username,
		&comment.Content,
		&comment.Likes,
		&comment.Dislikes,
		&comment.Created,
	); errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrCommentNotFound
	} else if err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *CommentsRepositorySqlite) Update(input domain.UpdateCommentInput) error {
	query := "UPDATE comments SET content = ?, edited = CURRENT_TIMESTAMP WHERE id = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(input.Content, input.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *CommentsRepositorySqlite) Delete(input domain.DeleteCommentInput) error {
	query := "DELETE FROM comments WHERE id = ?"
	stmt, err := r.db.Prepare(query)
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
