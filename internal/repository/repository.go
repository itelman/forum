package repository

import (
	"forum/internal/repository/models"
	"net/url"
)

type Repository struct {
	Users interface {
		Insert(string, string, string) error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
		InsertOAuth(string, string, string, string) (int, error)
		OAuth(string, string, string) (int, error)
	}
	Posts interface {
		Insert(int, string, string) (int, error)
		Get(int) (*models.Post, error)
		Delete(int) error
		Update(int, string, string) error
		Latest() ([]*models.Post, error)
		Filter(int, url.Values, func(int, int, string) (bool, error), func(int, []string) (bool, error)) ([]*models.Post, error)
		UpdateReactions(int, func(int) (int, error), func(int) (int, error)) error
		Created(int) ([]*models.Posts_Comments, error)
		Reacted(int, func(int) ([]*models.PostReaction, error)) ([]*models.Posts_Comments, error)
		Commented(int, func(int) ([]*models.Comment, error)) ([]*models.Posts_Comments, error)
	}
	Comments interface {
		Insert(int, int, string) error
		Get(int) (*models.Comment, error)
		Delete(int) (int, error)
		Update(int, string) error
		Latest(int) ([]*models.Comment, error)
		UpdateReactions(int, func(int) (int, error), func(int) (int, error)) error
		GetDistinctCommentsByUser(int) ([]*models.Comment, error)
		GetByUserForPost(int, int) ([]*models.Comment, error)
		LatestIgnoreUser(int, int) ([]*models.Comment, error)
	}
	Categories interface {
		Get(int) (*models.Category, error)
		Latest() ([]*models.Category, error)
	}
	Post_Category interface {
		Insert(int, []string) error
		Get(int) ([]string, error)
		FilterByCategories(int, []string) (bool, error)
	}
	Post_Reactions interface {
		Insert(int, int, int) error
		Get(int, int) (int, error)
		Delete(int, int) error
		Likes(int) (int, error)
		Dislikes(int) (int, error)
		FilterByLiked(int, int, string) (bool, error)
		GetReactionsByUser(int) ([]*models.PostReaction, error)
		LatestIgnoreUser(int, int) ([]*models.PostReaction, error)
	}
	Comment_Reactions interface {
		Insert(int, int, int) error
		Get(int, int) (int, error)
		Delete(int, int) error
		Likes(int) (int, error)
		Dislikes(int) (int, error)
	}
	Images interface {
		Insert(int, string) error
		Get(int) (*models.Image, error)
	}
}
