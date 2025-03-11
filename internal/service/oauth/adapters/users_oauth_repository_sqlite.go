package adapters

import (
	"database/sql"
	"errors"
	"github.com/itelman/forum/internal/service/oauth/domain"
)

type UsersOAuthRepositorySqlite struct {
	db *sql.DB
}

func NewUsersOAuthRepositorySqlite(db *sql.DB) *UsersOAuthRepositorySqlite {
	return &UsersOAuthRepositorySqlite{db}
}

func (r *UsersOAuthRepositorySqlite) GetUserID(input domain.GetUserOAuthInput) (int, error) {
	query := "SELECT user_id FROM users_oauth WHERE oauth_type_id = ? AND account_id = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	var userID int
	if err := stmt.QueryRow(input.AuthTypeID, input.AccountID).Scan(&userID); errors.Is(err, sql.ErrNoRows) {
		return -1, domain.ErrOAuthUserNotFound
	} else if err != nil {
		return -1, err
	}

	return userID, nil
}

func (r *UsersOAuthRepositorySqlite) Create(input domain.CreateUserOAuthInput) error {
	query := "INSERT INTO users_oauth (user_id, oauth_type_id, account_id) VALUES (?, ?, ?)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(input.UserID, input.AuthTypeID, input.AccountID)
	if err != nil {
		return err
	}

	return nil
}
