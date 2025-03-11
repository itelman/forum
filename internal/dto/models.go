package dto

import (
	"time"
)

type User struct {
	ID       int
	Username string
	Email    string
	Created  time.Time
}

type Post struct {
	ID               int
	User             *User
	Title            string
	Content          string
	Categories       []string
	Image            *Image
	Comments         []*Comment
	Likes            int
	Dislikes         int
	Created          time.Time
	AuthUserReaction int
}

type Comment struct {
	ID               int
	PostID           int
	User             *User
	Content          string
	Likes            int
	Dislikes         int
	Created          time.Time
	AuthUserReaction int
}

type PostReaction struct {
	ID      int
	PostID  int
	User    *User
	IsLike  int
	Created time.Time
}

type CommentReaction struct {
	ID        int
	CommentID int
	User      *User
	IsLike    int
	Created   time.Time
}

type Image struct {
	ID       int
	PostID   int
	Path     string
	Uploaded time.Time
}

type Category struct {
	ID      int
	Name    string
	Created time.Time
}

type Request struct {
	ID      int
	User    *User
	Created time.Time
}

type Report struct {
	ID          int
	Post        *Post
	Moderator   *User
	Content     string
	AdminReview string
	Created     time.Time
	Reviewed    time.Time
}
