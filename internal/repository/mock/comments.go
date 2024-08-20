package mock

import (
	"forum/internal/repository/models"
)

type CommentModel struct {
	DB []*models.Comment
}

func NewCommentModel() *CommentModel {
	return &CommentModel{make([]*models.Comment, 0)}
}

func (m *CommentModel) Insert(post_id, user_id, content string) error {
	return nil
}

func (m *CommentModel) Latest(post_id int) ([]*models.Comment, error) {
	comments := []*models.Comment{}

	return comments, nil
}

func (m *CommentModel) UpdateReactions(id int, Likes func(int) (int, error), Dislikes func(int) (int, error)) error {
	return nil
}
