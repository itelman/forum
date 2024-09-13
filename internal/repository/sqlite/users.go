package sqlite

import (
	"database/sql"
	"forum/internal/repository/models"
	"strings"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{db}
}

func (m *UserModel) Insert(name, email, password string) error {
	name = strings.ToLower(name)
	password = strings.ToLower(password)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password) VALUES (?, ?, ?)`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return models.ErrDuplicateNameOrEmail
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(name, password string) (int, error) {
	name = strings.ToLower(name)

	var id int
	var hashedPassword []byte
	var password_db sql.NullString

	row := m.DB.QueryRow("SELECT id, hashed_password FROM users WHERE name=?", name)
	err := row.Scan(&id, &password_db)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	if password_db.Valid {
		hashedPassword = []byte(password_db.String)
	} else {
		return 0, models.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	var email_db sql.NullString

	s := &models.User{}
	stmt := `SELECT id, name, email, created FROM users WHERE id=?`
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Name, &email_db, &s.Created)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	if email_db.Valid {
		s.Email = email_db.String
	}

	return s, nil
}
