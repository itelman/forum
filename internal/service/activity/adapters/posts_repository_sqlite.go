package adapters

import (
	"database/sql"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/activity/domain"
)

type PostsRepositorySqlite struct {
	db *sql.DB
}

func NewPostsRepositorySqlite(db *sql.DB) *PostsRepositorySqlite {
	return &PostsRepositorySqlite{db}
}

func (r *PostsRepositorySqlite) GetAllCreated(input domain.GetAllCreatedPostsInput) ([]*dto.Post, error) {
	query := "SELECT posts.id, users.username, posts.title, posts.content, posts.likes, posts.dislikes, posts.created, COALESCE(post_reactions.is_like, -1) AS is_like FROM posts INNER JOIN users ON posts.user_id = users.id LEFT JOIN post_reactions ON posts.id = post_reactions.post_id AND (post_reactions.user_id = ?) WHERE posts.user_id = ?"
	if input.SortedByNewest {
		query += " ORDER BY posts.created DESC"
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

	posts := []*dto.Post{}
	for rows.Next() {
		post := &dto.Post{User: &dto.User{}}

		if err := rows.Scan(
			&post.ID,
			&post.User.Username,
			&post.Title,
			&post.Content,
			&post.Likes,
			&post.Dislikes,
			&post.Created,
			&post.AuthUserReaction,
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

func (r *PostsRepositorySqlite) GetAllReacted(input domain.GetAllReactedPostsInput) ([]*dto.Post, error) {
	query := "SELECT p.id, u.username, p.title, p.content, p.likes, p.dislikes, p.created, pr.is_like FROM posts p INNER JOIN users u ON p.user_id = u.id JOIN post_reactions pr ON p.id = pr.post_id WHERE pr.user_id = ?"
	if input.SortedByNewest {
		query += " ORDER BY p.created DESC"
	}
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(input.AuthUserID)
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
			&post.Content,
			&post.Likes,
			&post.Dislikes,
			&post.Created,
			&post.AuthUserReaction,
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

func (r *PostsRepositorySqlite) GetAllCommented(input domain.GetAllCommentedPostsInput) ([]*dto.Post, error) {
	query := "SELECT DISTINCT p.id, u.username, p.title, p.content, p.likes, p.dislikes, p.created, COALESCE(pr.is_like, -1) AS is_like FROM posts p INNER JOIN users u ON p.user_id = u.id JOIN comments c ON p.id = c.post_id LEFT JOIN post_reactions pr ON p.id = pr.post_id AND (pr.user_id = ?) WHERE c.user_id = ?"
	if input.SortedByNewest {
		query += " ORDER BY p.created DESC"
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

	posts := []*dto.Post{}
	for rows.Next() {
		post := &dto.Post{User: &dto.User{}}

		if err := rows.Scan(
			&post.ID,
			&post.User.Username,
			&post.Title,
			&post.Content,
			&post.Likes,
			&post.Dislikes,
			&post.Created,
			&post.AuthUserReaction,
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
