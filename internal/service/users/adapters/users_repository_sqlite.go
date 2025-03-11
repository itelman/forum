package adapters

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/users/domain"
	"golang.org/x/crypto/bcrypt"
)

type UsersRepositorySqlite struct {
	db *sql.DB
}

func NewUsersRepositorySqlite(db *sql.DB) *UsersRepositorySqlite {
	return &UsersRepositorySqlite{db}
}

func (r *UsersRepositorySqlite) Create(input domain.RegisterUserInput) error {
	pwdHashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		return err
	}

	query := "INSERT INTO users (username, email, hashed_password) VALUES (?, ?, ?)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(input.Username, input.Email, string(pwdHashed))
	if err != nil {
		return err
	}

	return nil
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
		return nil, domain.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	if emailSql.Valid {
		user.Email = emailSql.String
	}

	return user, nil
}

func (r *UsersRepositorySqlite) Authenticate(input domain.AuthUserInput) (int, error) {
	var id int
	var pwd_hashed []byte
	var pwd_db sql.NullString

	query := "SELECT id, hashed_password FROM users WHERE username = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	if err := stmt.QueryRow(input.Username).Scan(&id, &pwd_db); errors.Is(err, sql.ErrNoRows) {
		return -1, domain.ErrUserNotFound
	} else if err != nil {
		return -1, err
	}

	if pwd_db.Valid {
		pwd_hashed = []byte(pwd_db.String)
	} else {
		return -1, domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(pwd_hashed, []byte(input.Password)); errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return -1, domain.ErrInvalidCredentials
	} else if err != nil {
		return -1, err
	}

	return id, nil
}
