package sqlite

import (
	"database/sql"
	"forum/internal/repository"
)

func NewRepository(db *sql.DB) *repository.Repository {
	return &repository.Repository{
		Users:             NewUserModel(db),
		Posts:             NewPostModel(db),
		Comments:          NewCommentModel(db),
		Categories:        NewCategoryModel(db),
		Post_Category:     NewPostCategoryModel(db),
		Post_Reactions:    NewPostReactionModel(db),
		Comment_Reactions: NewCommentReactionModel(db),
		Images:            NewImageModel(db),
	}
}
