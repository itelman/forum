package sqlite

import (
	"database/sql"
	"errors"
	"forum/pkg/models"
	"net/url"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(user_id, title, content, likes, dislikes string) (int, error) {
	var stmt string

	stmt = `INSERT INTO snippets (user_id, title, content, likes, dislikes)
		VALUES(?, ?, ?, ?, ?)`
	result, err := m.DB.Exec(stmt, user_id, title, content, likes, dislikes)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	// stmt := `SELECT id, title, content, created, likes, dislikes FROM snippets WHERE id = ?`
	stmt := `SELECT snippets.id, users.name, snippets.title, snippets.content, snippets.created, snippets.likes, snippets.dislikes FROM snippets INNER JOIN users ON snippets.user_id = users.id WHERE snippets.id = ?`
	row := m.DB.QueryRow(stmt, id)

	s := &models.Snippet{}
	err := row.Scan(&s.ID, &s.Username, &s.Title, &s.Content, &s.Created, &s.Likes, &s.Dislikes)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT snippets.id, users.id, users.name, snippets.title, snippets.content, snippets.created, snippets.likes, snippets.dislikes FROM snippets INNER JOIN users ON snippets.user_id = users.id ORDER BY snippets.created DESC`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		s := &models.Snippet{}
		err := rows.Scan(&s.ID, &s.UserID, &s.Username, &s.Title, &s.Content, &s.Created, &s.Likes, &s.Dislikes)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}

func (m *SnippetModel) UpdateReactions(id int, Likes func(int) (int, error), Dislikes func(int) (int, error)) error {
	likes, err := Likes(id)
	if err != nil {
		return err
	}

	dislikes, err := Dislikes(id)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec("UPDATE snippets SET likes = $1 WHERE id = $2", likes, id)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec("UPDATE snippets SET dislikes = $1 WHERE id = $2", dislikes, id)
	if err != nil {
		return err
	}

	return nil
}

/*func (m *SnippetModel) UpdateReactions(id, userID, content, likes, dislikes string) (int, int, error) {
	var err error
	var updatedLikes, updatedDislikes int

	switch likes {
	case "like":
		_, err = m.DB.Exec("UPDATE snippets SET likes = likes + 1 WHERE id = $1", id)
		if err != nil {
			return 0, 0, err
		}
	case "dislike":
		_, err = m.DB.Exec("UPDATE snippets SET dislikes = dislikes + 1 WHERE id = $1", id)
		if err != nil {
			return 0, 0, err
		}
	}

	// Retrieve updated counts after the update
	err = m.DB.QueryRow("SELECT likes, dislikes FROM snippets WHERE id = $1", id).Scan(&updatedLikes, &updatedDislikes)
	if err != nil {
		return 0, 0, err
	}

	return updatedLikes, updatedDislikes, nil
}*/

func (m *SnippetModel) Filter(form url.Values, FilterByLiked func(int, string, string) (bool, error),
	FilterByCategories func(int, []string, int) (bool, error)) ([]*models.Snippet, error) {
	var results []*models.Snippet

	snippets, err := m.Latest()
	if err != nil {
		return nil, err
	}

	for _, post := range snippets {
		cond_created := m.FilterByCreated(post.UserID, form.Get("user_id"), form.Get("created"))

		cond_liked, err := FilterByLiked(post.ID, form.Get("user_id"), form.Get("liked"))
		if err != nil {
			return nil, err
		}

		cond_categories, err := FilterByCategories(post.ID, form["categories"], len(form.Get("categories")))
		if err != nil {
			return nil, err
		}

		if cond_created && cond_liked && cond_categories {
			results = append(results, post)
		}
	}

	return results, nil
}

func (m *SnippetModel) FilterByCreated(post_user, user_id, val string) bool {
	if val != "1" {
		return true
	}

	if post_user == user_id {
		return true
	}

	return false
}

func (m *SnippetModel) Paginate(snippets []*models.Snippet, page, snippetNum int) ([]*models.Snippet, int, error) {
	if len(snippets) == 0 {
		return nil, 0, nil
	}

	pages := len(snippets) / snippetNum

	if len(snippets)%snippetNum != 0 {
		pages++
	}

	if page > pages {
		return nil, -1, errors.New("404")
	}

	start := (page - 1) * snippetNum
	end := page * snippetNum

	if end > len(snippets) {
		end = len(snippets)
	}

	return snippets[start:end], pages, nil
}
