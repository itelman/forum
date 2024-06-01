package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          int
	Username    string
	Password    string
	Email       string
	CreatedTime time.Time
}

func (s *Storage) GetUserBy(arg string, value interface{}) (User, error) {
	var user User
	query := fmt.Sprintf("SELECT id, username, password, email, created_time FROM users WHERE %s = $1", arg)
	err := s.db.QueryRow(query, value).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedTime)

	return user, err
}

func (s *Storage) CheckCredentials(user User) (bool, error) {
	var passwordDB string
	err := s.db.QueryRow("SELECT password FROM users WHERE username = $1", user.Username).Scan(&passwordDB)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	if !user.VerifyPassword(passwordDB) {
		return false, nil
	}

	return true, nil
}

func (s *Storage) InsertUser(user User) error {
	isEmailUnique, err := s.IsEmailUnique(user)
	if err != nil {
		return err
	}

	isUsernameUnique, err := s.IsUsernameUnique(user)
	if err != nil {
		return err
	}

	if !(isEmailUnique && isUsernameUnique) {
		return errors.New("Email/username exists")
	}

	err = user.HashPassword()
	if err != nil {
		return err
	}

	_, err = s.db.Exec("INSERT INTO users (username, password, email) VALUES ($1, $2, $3)", user.Username, user.Password, user.Email)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) ChangePassword(user User) error {
	err := user.HashPassword()
	if err != nil {
		return err
	}

	_, err = s.db.Exec("UPDATE users SET password = $1 WHERE id = $2", user.Password, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) ChangeEmail(user User) error {
	isEmailUnique, err := s.IsEmailUnique(user)
	if err != nil {
		return err
	}

	if !isEmailUnique {
		return errors.New("Email exists")
	}

	_, err = s.db.Exec("UPDATE users SET email = $1 WHERE id = $2", user.Email, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) IsUsernameUnique(user User) (bool, error) {
	_, err := s.GetUserBy("username", user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, err
	}

	return false, nil
}

func (s *Storage) IsEmailUnique(user User) (bool, error) {
	_, err := s.GetUserBy("email", user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, err
	}

	return false, nil
}

func (u *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return err
	}

	u.Password = string(bytes)

	return nil
}

func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(u.Password))

	return err == nil
}
