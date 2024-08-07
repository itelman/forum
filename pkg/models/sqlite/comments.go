package sqlite

import (
	"database/sql"
	"forum/pkg/models"
)

type CommentModel struct {
	DB *sql.DB
}

func (m *CommentModel) Insert(snippet_id, user_id, content string) error {
	var stmt string

	stmt = `INSERT INTO comments (snippet_id, user_id, content) VALUES ($1, $2, $3)`
	_, err := m.DB.Exec(stmt, snippet_id, user_id, content)
	if err != nil {
		return err
	}

	return err
}

func (m *CommentModel) Latest(snippet_id int) ([]*models.Comment, error) {
	stmt := `SELECT comments.id, comments.snippet_id, users.name, comments.content, comments.created, comments.likes, comments.dislikes FROM comments INNER JOIN users ON comments.user_id = users.id WHERE comments.snippet_id = ? ORDER BY comments.created DESC`
	rows, err := m.DB.Query(stmt, snippet_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*models.Comment{}

	for rows.Next() {
		s := &models.Comment{}
		err := rows.Scan(&s.ID, &s.SnippetID, &s.Username, &s.Content, &s.Created, &s.Likes, &s.Dislikes)
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
