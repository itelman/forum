package adapters

import (
	"database/sql"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/filters/domain"
)

type PostsRepositorySqlite struct {
	db *sql.DB
}

func NewPostsRepositorySqlite(db *sql.DB) *PostsRepositorySqlite {
	return &PostsRepositorySqlite{db}
}

func (r *PostsRepositorySqlite) GetManyByFilters(input domain.GetPostsByFiltersInput) ([]*dto.Post, error) {
	baseQuery := "SELECT posts.id, users.id, users.username, posts.title, posts.created FROM posts INNER JOIN users ON posts.user_id = users.id"

	catgClause := " WHERE EXISTS (SELECT 1 FROM post_categories WHERE post_categories.post_id = posts.id AND post_categories.category_id = ?)"
	createdClause := " posts.user_id = ?"
	likedClause := " EXISTS (SELECT 1 FROM post_reactions WHERE post_reactions.post_id = posts.id AND post_reactions.user_id = ? AND post_reactions.is_like = 1)"

	prevFilterApplied := false
	args := make([]interface{}, 0)

	if input.CategoryID != -1 {
		baseQuery += catgClause
		args = append(args, input.CategoryID)
		prevFilterApplied = true
	}

	if input.Created {
		if prevFilterApplied {
			baseQuery += " AND"
		} else {
			baseQuery += " WHERE"
		}
		args = append(args, input.AuthUserID)
		baseQuery += createdClause
		prevFilterApplied = true
	}

	if input.Liked {
		if prevFilterApplied {
			baseQuery += " AND"
		} else {
			baseQuery += " WHERE"
		}
		args = append(args, input.AuthUserID)
		baseQuery += likedClause
	}

	if input.SortedByNewest {
		baseQuery += " ORDER BY posts.created DESC"
	}

	stmt, err := r.db.Prepare(baseQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*dto.Post{}
	for rows.Next() {
		post := &dto.Post{User: &dto.User{}}

		if err := rows.Scan(
			&post.ID,
			&post.User.ID,
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
