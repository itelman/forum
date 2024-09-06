package sqlite

import (
	"database/sql"
	"fmt"
	"forum/internal/repository/models"
	"strings"

	"github.com/mattn/go-sqlite3"
)

func (m *UserModel) InsertOAuth(name, email string, provider_id string, providerName string) (int, error) {
	name = strings.ToLower(name)
	email = strings.ToLower(email)

	var emailNull sql.NullString

	if email == "" {
		emailNull = sql.NullString{String: "", Valid: false}
	} else {
		emailNull = sql.NullString{String: email, Valid: true}
	}

	stmt := fmt.Sprintf("INSERT INTO users (name, email, %s) VALUES (?, ?, ?)", providerName)
	result, err := m.DB.Exec(stmt, name, emailNull, provider_id)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return -1, models.ErrDuplicateNameOrEmail
		}
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

func (m *UserModel) OAuth(email string, provider_id string, providerName string) (int, error) {
	var id int
	var provider_id_db sql.NullString

	stmt := fmt.Sprintf("SELECT id, %s FROM users WHERE %s = $1 OR email = $2", providerName, providerName)
	err := m.DB.QueryRow(stmt, provider_id, email).Scan(&id, &provider_id_db)
	if err == sql.ErrNoRows {
		return -1, models.ErrNoRecord
	} else if err != nil {
		return -1, err
	}

	stmt = fmt.Sprintf(`UPDATE users SET %s = $1 WHERE id = $2`, providerName)
	_, err = m.DB.Exec(stmt, provider_id, id)
	if err != nil {
		return -1, err
	}

	return id, nil
}
