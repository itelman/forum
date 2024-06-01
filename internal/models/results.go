package models

import (
	"net/url"
	"strconv"
)

type Results struct {
	Posts      []*Post
	Categories []Category
}

func GetPostsResults(form url.Values, db *Storage) ([]*Post, error) {
	posts, err := db.GetAllPosts()
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		post.SetReaction(db, 1)
	}

	if _, ok := form["category"]; ok {
		posts = FilterByCategories(posts, form["category"])
	}

	if _, ok := form["created"]; ok {
		posts = FilterByCreated(posts, 1)
	}

	if _, ok := form["liked"]; ok {
		posts = FilterByLiked(posts)
	}

	return posts, nil
}

func FilterByCategories(posts []*Post, categories_id []string) []*Post {
	var results []*Post
	catIDmap := make(map[string]struct{})

	for _, id := range categories_id {
		catIDmap[id] = struct{}{}
	}

	for _, post := range posts {
		for k, _ := range post.Categories {
			if _, ok := catIDmap[strconv.Itoa(k)]; ok {
				results = append(results, post)
				break
			}
		}
	}

	return results
}

func FilterByCreated(posts []*Post, user_id int) []*Post {
	var results []*Post

	for _, post := range posts {
		if post.UserID == user_id {
			results = append(results, post)
		}
	}

	return results
}

func FilterByLiked(posts []*Post) []*Post {
	var results []*Post

	for _, post := range posts {
		if post.UserReaction == 1 {
			results = append(results, post)
		}
	}

	return results
}

func (s *Storage) SetResults(posts []*Post) (Results, error) {
	categories, err := s.GetAllCategories()
	if err != nil {
		return Results{}, err
	}

	return Results{posts, categories}, nil
}
