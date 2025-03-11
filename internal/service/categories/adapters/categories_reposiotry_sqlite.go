package adapters

import (
	"database/sql"
	"errors"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/categories/domain"
)

type CategoriesRepositorySqlite struct {
	db *sql.DB
}

func NewCategoriesRepositorySqlite(db *sql.DB) *CategoriesRepositorySqlite {
	return &CategoriesRepositorySqlite{db}
}

func (r *CategoriesRepositorySqlite) Create(input domain.CreateCategoryInput) error {
	query := "INSERT INTO categories (name) VALUES(?)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(input.Name); err != nil {
		return err
	}

	return nil
}

func (r *CategoriesRepositorySqlite) Get(input domain.GetCategoryInput) (*dto.Category, error) {
	query := "SELECT id, name, created FROM categories WHERE id = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	category := &dto.Category{}
	if err := stmt.QueryRow(input.ID).Scan(
		&category.ID,
		&category.Name,
		&category.Created,
	); errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrCategoryNotFound
	} else if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *CategoriesRepositorySqlite) GetAll(input domain.GetAllCategoriesInput) ([]*dto.Category, error) {
	query := "SELECT id, name, created FROM categories"
	if input.SortedByNewest {
		query += " ORDER BY categories.created DESC"
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

	categories := []*dto.Category{}
	for rows.Next() {
		category := &dto.Category{}

		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Created,
		); err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoriesRepositorySqlite) Delete(input domain.DeleteCategoryInput) error {
	query := "DELETE FROM categories WHERE id = ?"
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
