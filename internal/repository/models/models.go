package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord             = errors.New("models: no matching record found")
	ErrInvalidCredentials   = errors.New("models: invalid credentials")
	ErrDuplicateNameOrEmail = errors.New("UNIQUE constraint failed: users.name or users.email")
)

type Post struct {
	ID            int
	UserID        string
	Username      string
	Title         string
	Content       string
	Created       time.Time
	Likes         int
	Dislikes      int
	ReactedByUser string
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type Comment struct {
	ID            int
	UserID        string
	PostID        int
	Username      string
	Content       string
	Created       time.Time
	Likes         int
	Dislikes      int
	ReactedByUser string
}

type Category struct {
	ID      int
	Name    string
	Created time.Time
}

type Error struct {
	Code    int
	Message string
}

type Session struct {
	Id      string
	User_Id int
}
