package adapters

import (
	"database/sql"
	"errors"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/posts/domain"
)

type PostsRepositorySqlite struct {
	db *sql.DB
}

func NewPostsRepositorySqlite(db *sql.DB) *PostsRepositorySqlite {
	return &PostsRepositorySqlite{db}
}

func (r *PostsRepositorySqlite) Create(tx *sql.Tx, input domain.CreatePostInput) (int, error) {
	query := "INSERT INTO posts (user_id, title, content) VALUES(?, ?, ?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(input.UserID, input.Title, input.Content)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

func (r *PostsRepositorySqlite) Get(input domain.GetPostInput) (*dto.Post, error) {
	query := "SELECT posts.id, users.id, users.username, posts.title, posts.content, posts.likes, posts.dislikes, posts.created, COALESCE(post_reactions.is_like, -1) AS is_like FROM posts INNER JOIN users ON posts.user_id = users.id LEFT JOIN post_reactions ON posts.id = post_reactions.post_id AND (post_reactions.user_id = ? OR ? = -1) WHERE posts.id = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	post := &dto.Post{User: &dto.User{}}
	if err := stmt.QueryRow(input.AuthUserID, input.AuthUserID, input.ID).Scan(
		&post.ID,
		&post.User.ID,
		&post.User.Username,
		&post.Title,
		&post.Content,
		&post.Likes,
		&post.Dislikes,
		&post.Created,
		&post.AuthUserReaction,
	); errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrPostNotFound
	} else if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *PostsRepositorySqlite) GetAll(input domain.GetAllPostsInput) ([]*dto.Post, error) {
	query := "SELECT posts.id, users.username, posts.title, posts.created FROM posts INNER JOIN users ON posts.user_id = users.id"
	if input.SortedByNewest {
		query += " ORDER BY posts.created DESC"
	}
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*dto.Post{}
	for rows.Next() {
		post := &dto.Post{User: &dto.User{}}

		if err := rows.Scan(
			&post.ID,
			&post.User.Username,
			&post.Title,
			&post.Created,
		); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostsRepositorySqlite) Update(input domain.UpdatePostInput) error {
	query := "UPDATE posts SET title = ?, content = ?, edited = CURRENT_TIMESTAMP WHERE id = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(input.Title, input.Content, input.ID); err != nil {
		return err
	}

	return nil
}

func (r *PostsRepositorySqlite) Delete(input domain.DeletePostInput) error {
	query := "DELETE FROM posts WHERE id = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(input.ID); err != nil {
		return err
	}

	return nil
}
