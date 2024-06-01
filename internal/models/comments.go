package models

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID         int
	PostID     int
	UserID     int
	Content    string
	PostedTime time.Time
	Likes      int
	Dislikes   int

	Username     string
	UserReaction int
}

func (s *Storage) InsertComment(comment Comment) error {
	_, err := s.db.Exec("INSERT INTO comments (post_id, user_id, content) VALUES ($1, $2, $3)", comment.PostID, comment.UserID, comment.Content)

	return err
}

func (s *Storage) GetAllCommentsForPost(post_id int) ([]*Comment, error) {
	rows, err := s.db.Query("SELECT id, post_id, user_id, content, posted_time, likes, dislikes FROM comments WHERE post_id = $1", post_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.PostedTime, &comment.Likes, &comment.Dislikes)
		if err != nil {
			return nil, err
		}

		err = comment.SetUsername(s)
		if err != nil {
			return nil, err
		}

		err = comment.SetTimezone("Asia/Oral")
		if err != nil {
			return nil, err
		}

		comments = append(comments, &comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (c *Comment) SetUsername(s *Storage) error {
	user_id := c.UserID

	user, err := s.GetUserBy("id", user_id)
	if err != nil {
		return err
	}
	c.Username = user.Username

	return nil
}

func (c *Comment) SetTimezone(location string) error {
	utc5, err := time.LoadLocation(location)
	if err != nil {
		return err
	}
	c.PostedTime = c.PostedTime.In(utc5)

	return nil
}

func (c *Comment) SetReaction(s *Storage, user_id int) error {
	reaction, err := s.GetCommentReaction(c.ID, user_id)
	if err != nil && !(err == sql.ErrNoRows) {
		return err
	}

	c.UserReaction = reaction.LikeOrDislike
	return nil
}

func (s *Storage) UpdateCommentReactions(comment_id int) error {
	/*switch reaction {
	case -1:
		_, err := s.db.Exec("UPDATE comments SET likes = likes - 1 WHERE id = $1", comment_id)
		return err
	case -2:
		_, err := s.db.Exec("UPDATE comments SET dislikes = dislikes - 1 WHERE id = $1", comment_id)
		return err
	case 1:
		_, err := s.db.Exec("UPDATE comments SET likes = likes + 1 WHERE id = $1", comment_id)
		return err
	case 2:
		_, err := s.db.Exec("UPDATE comments SET dislikes = dislikes + 1 WHERE id = $1", comment_id)
		return err
	}*/

	likes, err := s.CountCommentLikes(comment_id)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("UPDATE comments SET likes = $1 WHERE id = $2", likes, comment_id)
	if err != nil {
		return err
	}

	dislikes, err := s.CountCommentDislikes(comment_id)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("UPDATE comments SET dislikes = $1 WHERE id = $2", dislikes, comment_id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CountCommentLikes(post_id int) (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM comment_reactions WHERE comment_id = $1 AND like_or_dislike = $2", post_id, 1).Scan(&count)
	if err != nil && !(err == sql.ErrNoRows) {
		return -1, err
	}

	return count, nil
}

func (s *Storage) CountCommentDislikes(post_id int) (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM comment_reactions WHERE comment_id = $1 AND like_or_dislike = $2", post_id, 2).Scan(&count)
	if err != nil && !(err == sql.ErrNoRows) {
		return -1, err
	}

	return count, nil
}
