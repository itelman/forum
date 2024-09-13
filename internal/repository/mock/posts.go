package mock

import (
	"forum/internal/repository/models"
	"net/url"
)

type PostModel struct {
	DB []*models.Post
}

func NewPostModel() *PostModel {
	return &PostModel{make([]*models.Post, 0)}
}

func (m *PostModel) Insert(user_id int, title, content string) (int, error) {
	id := len(m.DB) + 1

	m.DB = append(m.DB, &models.Post{ID: id, UserID: user_id, Title: title, Content: content})

	return id, nil
}

func (m *PostModel) Get(id int) (*models.Post, error) {
	if id < 1 || id > len(m.DB) {
		return &models.Post{}, models.ErrNoRecord
	}

	post := m.DB[id-1]
	return post, nil
}

func (m *PostModel) Delete(id int) error {
	return nil
}

func (m *PostModel) Update(id int, title, content string) error {
	return nil
}

func (m *PostModel) Latest() ([]*models.Post, error) {
	posts := []*models.Post{}

	return posts, nil
}

func (m *PostModel) UpdateReactions(id int, Likes func(int) (int, error), Dislikes func(int) (int, error)) error {
	return nil
}

func (m *PostModel) Filter(user_id int, form url.Values, FilterByLiked func(int, int, string) (bool, error),
	FilterByCategories func(int, []string) (bool, error)) ([]*models.Post, error) {
	var results []*models.Post

	return results, nil
}

func (m *PostModel) Created(id int) ([]*models.Posts_Comments, error) {
	return nil, nil
}

func (m *PostModel) Reacted(id int, GetReactionsByUser func(int) ([]*models.PostReaction, error)) ([]*models.Posts_Comments, error) {
	return nil, nil
}

func (m *PostModel) Commented(id int, GetDistinctCommentsByUser func(int) ([]*models.Comment, error)) ([]*models.Posts_Comments, error) {
	return nil, nil
}
