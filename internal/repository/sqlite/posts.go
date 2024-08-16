package sqlite

import (
	"database/sql"
	"forum/internal/repository/models"
	"net/url"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(user_id, title, content string) (int, error) {
	stmt := `INSERT INTO posts (user_id, title, content) VALUES(?, ?, ?)`
	result, err := m.DB.Exec(stmt, user_id, title, content)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

func (m *PostModel) Get(id int) (*models.Post, error) {
	stmt := `SELECT posts.id, users.name, posts.title, posts.content, posts.created, posts.likes, posts.dislikes FROM posts INNER JOIN users ON posts.user_id = users.id WHERE posts.id = ?`
	row := m.DB.QueryRow(stmt, id)

	s := &models.Post{}
	err := row.Scan(&s.ID, &s.Username, &s.Title, &s.Content, &s.Created, &s.Likes, &s.Dislikes)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

func (m *PostModel) Latest() ([]*models.Post, error) {
	stmt := `SELECT posts.id, users.id, users.name, posts.title, posts.content, posts.created, posts.likes, posts.dislikes FROM posts INNER JOIN users ON posts.user_id = users.id ORDER BY posts.created DESC`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*models.Post{}

	for rows.Next() {
		s := &models.Post{}
		var likes sql.NullInt64
		var dislikes sql.NullInt64
		if likes.Valid {
			s.Likes = int(likes.Int64)
		} else {
			s.Likes = 0
		}

		if dislikes.Valid {
			s.Dislikes = int(dislikes.Int64)
		} else {
			s.Dislikes = 0
		}
		err := rows.Scan(&s.ID, &s.UserID, &s.Username, &s.Title, &s.Content, &s.Created, &s.Likes, &s.Dislikes)
		if err != nil {
			return nil, err
		}
		posts = append(posts, s)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (m *PostModel) UpdateReactions(id int, Likes func(int) (int, error), Dislikes func(int) (int, error)) error {
	l, err := Likes(id)
	if err != nil {
		return err
	}

	d, err := Dislikes(id)
	if err != nil {
		return err
	}
	_, err = m.DB.Exec("UPDATE posts SET likes = $1 WHERE id = $2", l, id)
	if err != nil {
		return err
	}
	_, err = m.DB.Exec("UPDATE posts SET dislikes = $1 WHERE id = $2", d, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostModel) Filter(user_id int, form url.Values, FilterByLiked func(int, int, string) (bool, error),
	FilterByCategories func(int, []string) (bool, error)) ([]*models.Post, error) {
	var results []*models.Post

	posts, err := m.Latest()
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		cond_created := m.FilterByCreated(post.UserID, user_id, form.Get("created"))

		cond_liked, err := FilterByLiked(post.ID, user_id, form.Get("liked"))
		if err != nil {
			return nil, err
		}

		cond_categories, err := FilterByCategories(post.ID, form["categories"])
		if err != nil {
			return nil, err
		}

		if cond_created && cond_liked && cond_categories {
			results = append(results, post)
		}
	}

	return results, nil
}

func (m *PostModel) FilterByCreated(post_user_id, user_id int, val string) bool {
	if len(val) == 0 || user_id == -1 {
		return true
	}

	if post_user_id == user_id {
		return true
	}

	return false
}
