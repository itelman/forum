package adapters

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/oauth/domain"
)

type UsersRepositorySqlite struct {
	db *sql.DB
}

func NewUsersRepositorySqlite(db *sql.DB) *UsersRepositorySqlite {
	return &UsersRepositorySqlite{db}
}

func (r *UsersRepositorySqlite) Get(input domain.GetUserInput) (*dto.User, error) {
	query := fmt.Sprintf("SELECT id, username, email, created FROM users WHERE %s = ?", input.Key)
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	user := &dto.User{}
	var emailSql sql.NullString
	if err := stmt.QueryRow(input.Value).Scan(
		&user.ID,
		&user.Username,
		&emailSql,
		&user.Created,
	); errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrOAuthUserNotFound
	} else if err != nil {
		return nil, err
	}

	if emailSql.Valid {
		user.Email = emailSql.String
	}

	return user, nil
}
