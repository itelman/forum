package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateName      = errors.New("UNIQUE constraint failed: users.name")
	ErrDuplicateEmail     = errors.New("UNIQUE constraint failed: users.email")
)

type Snippet struct {
	ID       int
	UserID   string
	Username string
	Title    string
	Content  string
	Created  time.Time
	Likes    int
	Dislikes int
}

type User struct {
	ID             int
	Name           string
	Phone          string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type Comment struct {
	ID        int
	SnippetID int
	Username  string
	Content   string
	Created   time.Time
	Likes     int
	Dislikes  int
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

/*type PostCategory struct {
	PostID       int
	CategoryName string
	Created      time.Time
}

type PostReaction struct {
	ID      int
	PostID  int
	UserID  int
	IsLike  int
	Created time.Time
}

type CommentReaction struct {
	ID        int
	CommentID int
	UserID    int
	IsLike    int
	Created   time.Time
}*/
