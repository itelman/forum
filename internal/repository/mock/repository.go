package mock

import "forum/internal/repository"

func NewRepository() *repository.Repository {
	return &repository.Repository{
		Posts:         NewPostModel(),
		Comments:      NewCommentModel(),
		Post_Category: NewPostCategoryModel(),
		Images:        NewImageModel(),
	}
}
