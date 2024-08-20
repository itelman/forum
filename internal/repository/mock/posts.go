package mock

import (
	"forum/internal/repository/models"
	"net/url"
	"strconv"
)

type PostModel struct {
	DB []*models.Post
}

func NewPostModel() *PostModel {
	return &PostModel{make([]*models.Post, 0)}
}

func (m *PostModel) Insert(user_id, title, content string) (int, error) {
	id := len(m.DB) + 1
	userID, err := strconv.Atoi(user_id)
	if err != nil {
		return -1, err
	}

	m.DB = append(m.DB, &models.Post{ID: id, UserID: userID, Title: title, Content: content})

	return id, nil
}

func (m *PostModel) Get(id int) (*models.Post, error) {
	if id < 1 || id > len(m.DB) {
		return &models.Post{}, models.ErrNoRecord
	}

	post := m.DB[id-1]
	return post, nil
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

func (m *PostModel) FilterByCreated(post_user_id, user_id int, val string) bool {
	return false
}
