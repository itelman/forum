package adapters

import (
	"database/sql"
	"errors"
	"github.com/itelman/forum/internal/service/posts/domain"
	"github.com/mattn/go-sqlite3"
)

type PostCategoriesRepositorySqlite struct {
	db *sql.DB
}

func NewPostCategoriesRepositorySqlite(db *sql.DB) *PostCategoriesRepositorySqlite {
	return &PostCategoriesRepositorySqlite{db}
}

func (r *PostCategoriesRepositorySqlite) Create(tx *sql.Tx, input domain.CreatePostCategoriesInput) error {
	query := "INSERT INTO post_categories (post_id, category_id) VALUES(?, ?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, catgId := range input.CategoriesID {
		_, err := stmt.Exec(input.PostID, catgId)

		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintForeignKey) {
			return domain.ErrPostsBadRequest
		} else if err != nil {
			return err
		}
	}

	return nil
}

func (r *PostCategoriesRepositorySqlite) GetAllForPost(input domain.GetPostCategoriesInput) ([]string, error) {
	query := "SELECT categories.name FROM post_categories INNER JOIN categories ON post_categories.category_id = categories.id WHERE post_categories.post_id = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(input.PostID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
