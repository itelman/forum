package models

import (
	"database/sql"
	"strconv"
	"time"
)

type Post struct {
	ID         int
	UserID     int
	Title      string
	Content    string
	PostedTime time.Time
	Likes      int
	Dislikes   int

	Username     string
	Comments     []*Comment
	UserReaction int
	Categories   map[int]Category
}

func (s *Storage) InsertPost(post Post, categories_id []string) error {
	_, err := s.db.Exec("INSERT INTO posts (id, user_id, title, content) VALUES ($1, $2, $3, $4)", post.ID, post.UserID, post.Title, post.Content)
	if err != nil {
		return err
	}

	for _, id_str := range categories_id {
		id, err := strconv.Atoi(id_str)
		if err != nil {
			return err
		}

		err = s.InsertPostCategory(post.ID, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) GetPostByID(id int) (Post, error) {
	var post Post
	query := "SELECT id, user_id, title, content, posted_time, likes, dislikes FROM posts WHERE id = $1"
	err := s.db.QueryRow(query, id).Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.PostedTime, &post.Likes, &post.Dislikes)
	if err != nil {
		return Post{}, err
	}

	err = post.SetUsername(s)
	if err != nil {
		return Post{}, err
	}

	err = post.SetTimezone("Asia/Oral")
	if err != nil {
		return Post{}, err
	}

	comments, err := s.GetAllCommentsForPost(id)
	if err != nil {
		return Post{}, err
	}
	post.Comments = reverseArray(comments)

	categories, err := post.GetCategoriesForPost(s)
	if err != nil {
		return Post{}, err
	}
	post.Categories = categories

	return post, nil
}

func reverseArray(arr []*Comment) []*Comment {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}

	return arr
}

func (s *Storage) UpdatePostReactions(post_id int) error {
	/*switch reaction {
	case -1:
		_, err := s.db.Exec("UPDATE posts SET likes = likes - 1 WHERE id = $1", post_id)
		return err
	case -2:
		_, err := s.db.Exec("UPDATE posts SET dislikes = dislikes - 1 WHERE id = $1", post_id)
		return err
	case 1:
		_, err := s.db.Exec("UPDATE posts SET likes = likes + 1 WHERE id = $1", post_id)
		return err
	case 2:
		_, err := s.db.Exec("UPDATE posts SET dislikes = dislikes + 1 WHERE id = $1", post_id)
		return err
	}*/

	likes, err := s.CountPostLikes(post_id)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("UPDATE posts SET likes = $1 WHERE id = $2", likes, post_id)
	if err != nil {
		return err
	}

	dislikes, err := s.CountPostDislikes(post_id)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("UPDATE posts SET dislikes = $1 WHERE id = $2", dislikes, post_id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CountPostLikes(post_id int) (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM post_reactions WHERE post_id = $1 AND like_or_dislike = $2", post_id, 1).Scan(&count)
	if err != nil && !(err == sql.ErrNoRows) {
		return -1, err
	}

	return count, nil
}

func (s *Storage) CountPostDislikes(post_id int) (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM post_reactions WHERE post_id = $1 AND like_or_dislike = $2", post_id, 2).Scan(&count)
	if err != nil && !(err == sql.ErrNoRows) {
		return -1, err
	}

	return count, nil
}

func (s *Storage) SetPostID() (int, error) {
	var id int
	query := "SELECT t1.id + 1 AS START FROM posts AS t1 LEFT OUTER JOIN posts AS t2 ON t1.id + 1 = t2.id WHERE t2.id IS NULL"
	err := s.db.QueryRow(query).Scan(&id)
	if err == sql.ErrNoRows {
		return 1, nil
	}

	return id, err
}

func (s *Storage) GetAllPosts() ([]*Post, error) {
	rows, err := s.db.Query("SELECT * FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.PostedTime, &post.Likes, &post.Dislikes)
		if err != nil {
			return nil, err
		}

		err = post.SetUsername(s)
		if err != nil {
			return nil, err
		}

		err = post.SetTimezone("Asia/Oral")
		if err != nil {
			return nil, err
		}

		categories, err := post.GetCategoriesForPost(s)
		if err != nil {
			return nil, err
		}
		post.Categories = categories

		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *Post) SetUsername(s *Storage) error {
	user_id := p.UserID

	user, err := s.GetUserBy("id", user_id)
	if err != nil {
		return err
	}
	p.Username = user.Username

	return nil
}

func (p *Post) SetTimezone(location string) error {
	utc5, err := time.LoadLocation(location)
	if err != nil {
		return err
	}
	p.PostedTime = p.PostedTime.In(utc5)

	return nil
}

func (p *Post) SetReaction(s *Storage, user_id int) error {
	reaction, err := s.GetPostReaction(p.ID, user_id)
	if err != nil && !(err == sql.ErrNoRows) {
		return err
	}

	p.UserReaction = reaction.LikeOrDislike
	return nil
}

func (p *Post) GetCategoriesForPost(s *Storage) (map[int]Category, error) {
	rows, err := s.db.Query("SELECT category_id FROM post_categories WHERE post_id = $1", p.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make(map[int]Category)
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		category, err := s.GetCategoryByID(id)
		if err != nil {
			return nil, err
		}

		categories[category.ID] = category
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
