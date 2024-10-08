package models

import (
	"errors"
	"time"
)

var (
	ErrBadGateway           = errors.New("models: error encountered from another server")
	ErrNoRecord             = errors.New("models: no matching record found")
	ErrInvalidCredentials   = errors.New("models: invalid credentials")
	ErrDuplicateNameOrEmail = errors.New("UNIQUE constraint failed: users.name or users.email")
)

type Post struct {
	ID            int
	UserID        int
	Username      string
	Title         string
	Content       string
	Created       time.Time
	Likes         int
	Dislikes      int
	Edited        time.Time
	ReactedByUser int
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
	UserID        int
	PostID        int
	Username      string
	Content       string
	Created       time.Time
	Likes         int
	Dislikes      int
	Edited        time.Time
	ReactedByUser int
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

type PostCategory struct {
	PostID       int
	CategoryName string
	Created      time.Time
}

type PostReaction struct {
	ID       int
	PostID   int
	UserID   int
	Username string
	IsLike   int
	Created  time.Time
}

type Posts_Comments struct {
	Post     *Post
	Comments []*Comment
}

type Image struct {
	ID       int
	PostID   int
	Path     string
	Uploaded time.Time
}
