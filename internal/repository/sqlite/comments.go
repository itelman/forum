package sqlite

import (
	"database/sql"
	"forum/internal/repository/models"
)

type CommentModel struct {
	DB *sql.DB
}

func (m *CommentModel) Insert(post_id, user_id, content string) error {
	var stmt string

	stmt = `INSERT INTO comments (post_id, user_id, content) VALUES ($1, $2, $3)`
	_, err := m.DB.Exec(stmt, post_id, user_id, content)
	if err != nil {
		return err
	}

	return err
}

func (m *CommentModel) Latest(post_id int) ([]*models.Comment, error) {
	stmt := `SELECT comments.id, comments.post_id, users.name, comments.content, comments.created, comments.likes, comments.dislikes FROM comments INNER JOIN users ON comments.user_id = users.id WHERE comments.post_id = ? ORDER BY comments.created DESC`
	rows, err := m.DB.Query(stmt, post_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*models.Comment{}

	for rows.Next() {
		s := &models.Comment{}
		err := rows.Scan(&s.ID, &s.PostID, &s.Username, &s.Content, &s.Created, &s.Likes, &s.Dislikes)
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
