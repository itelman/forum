package sqlite

import (
	"database/sql"
	"forum/internal/repository/models"

	"github.com/mattn/go-sqlite3"
)

type CommentModel struct {
	DB *sql.DB
}

func NewCommentModel(db *sql.DB) *CommentModel {
	return &CommentModel{db}
}

func (m *CommentModel) Insert(post_id, user_id int, content string) error {
	stmt := `INSERT INTO comments (post_id, user_id, content) VALUES ($1, $2, $3)`
	_, err := m.DB.Exec(stmt, post_id, user_id, content)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint && sqliteErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
			return models.ErrNoRecord
		}
		return err
	}

	return err
}

func (m *CommentModel) Get(id int) (*models.Comment, error) {
	stmt := `SELECT comments.id, comments.user_id, comments.post_id, users.name, comments.content, comments.created, comments.likes, comments.dislikes FROM comments INNER JOIN users ON comments.user_id = users.id WHERE comments.id = ?`
	row := m.DB.QueryRow(stmt, id)

	s := &models.Comment{}
	err := row.Scan(&s.ID, &s.UserID, &s.PostID, &s.Username, &s.Content, &s.Created, &s.Likes, &s.Dislikes)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

func (m *CommentModel) Delete(id int) error {
	stmt := `DELETE FROM comments WHERE id = ?`

	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *CommentModel) Update(id int, content string) error {
	stmt := `UPDATE comments SET content = $1, edited = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := m.DB.Exec(stmt, content, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *CommentModel) Latest(post_id int) ([]*models.Comment, error) {
	stmt := `SELECT comments.id, users.id, comments.post_id, users.name, comments.content, comments.created, comments.likes, comments.dislikes FROM comments INNER JOIN users ON comments.user_id = users.id WHERE comments.post_id = ? ORDER BY comments.created DESC`
	rows, err := m.DB.Query(stmt, post_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*models.Comment{}

	for rows.Next() {
		s := &models.Comment{}
		err := rows.Scan(&s.ID, &s.UserID, &s.PostID, &s.Username, &s.Content, &s.Created, &s.Likes, &s.Dislikes)
		if err != nil {
			return nil, err
		}
		comments = append(comments, s)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (m *CommentModel) LatestIgnoreUser(post_id, user_id int) ([]*models.Comment, error) {
	stmt := `SELECT comments.id, users.id, comments.post_id, users.name, comments.content, comments.created, comments.likes, comments.dislikes FROM comments INNER JOIN users ON comments.user_id = users.id WHERE comments.post_id = ? AND comments.user_id <> ? ORDER BY comments.created DESC`
	rows, err := m.DB.Query(stmt, post_id, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*models.Comment{}

	for rows.Next() {
		s := &models.Comment{}
		err := rows.Scan(&s.ID, &s.UserID, &s.PostID, &s.Username, &s.Content, &s.Created, &s.Likes, &s.Dislikes)
		if err != nil {
			return nil, err
		}
		comments = append(comments, s)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (m *CommentModel) UpdateReactions(id int, Likes func(int) (int, error), Dislikes func(int) (int, error)) error {
	likes, err := Likes(id)
	if err != nil {
		return err
	}

	dislikes, err := Dislikes(id)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec("UPDATE comments SET likes = $1 WHERE id = $2", likes, id)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec("UPDATE comments SET dislikes = $1 WHERE id = $2", dislikes, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *CommentModel) GetDistinctCommentsByUser(user_id int) ([]*models.Comment, error) {
	stmt := `SELECT id, post_id, user_id, content, created, likes, dislikes FROM comments WHERE user_id = ? GROUP BY post_id ORDER BY created DESC`
	rows, err := m.DB.Query(stmt, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*models.Comment{}

	for rows.Next() {
		s := models.Comment{}
		err := rows.Scan(&s.ID, &s.PostID, &s.UserID, &s.Content, &s.Created, &s.Likes, &s.Dislikes)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &s)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (m *CommentModel) GetByUserForPost(post_id, user_id int) ([]*models.Comment, error) {
	stmt := `SELECT comments.id, users.id, comments.post_id, users.name, comments.content, comments.created, comments.likes, comments.dislikes FROM comments INNER JOIN users ON comments.user_id = users.id WHERE comments.post_id = ? AND comments.user_id = ? ORDER BY comments.created DESC`
	rows, err := m.DB.Query(stmt, post_id, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*models.Comment{}

	for rows.Next() {
		s := models.Comment{}
		err := rows.Scan(&s.ID, &s.UserID, &s.PostID, &s.Username, &s.Content, &s.Created, &s.Likes, &s.Dislikes)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}
