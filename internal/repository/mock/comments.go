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

func (m *CommentModel) Insert(post_id, user_id int, content string) error {
	return nil
}

func (m *CommentModel) Get(id int) (*models.Comment, error) {
	return &models.Comment{}, nil
}

func (m *CommentModel) Latest(post_id int) ([]*models.Comment, error) {
	comments := []*models.Comment{}

	return comments, nil
}

func (m *CommentModel) UpdateReactions(id int, Likes func(int) (int, error), Dislikes func(int) (int, error)) error {
	return nil
}

func (m *CommentModel) Delete(id int) error {
	return nil
}

func (m *CommentModel) Update(id int, content string) error {
	return nil
}

func (m *CommentModel) LatestIgnoreUser(post_id, user_id int) ([]*models.Comment, error) {
	return nil, nil
}

func (m *CommentModel) GetDistinctCommentsByUser(user_id int) ([]*models.Comment, error) {
	return nil, nil
}

func (m *CommentModel) GetByUserForPost(post_id, user_id int) ([]*models.Comment, error) {
	return nil, nil
}
