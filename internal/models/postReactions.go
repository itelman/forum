package models

import "database/sql"

type PostReaction struct {
	PostID        int
	UserID        int
	LikeOrDislike int
}

func (s *Storage) InsertPostReaction(reaction PostReaction) error {
	if reaction.LikeOrDislike == -1 || reaction.LikeOrDislike == -2 {
		err := s.DeletePostReaction(reaction)
		return err
	}

	_, err := s.GetPostReaction(reaction.PostID, reaction.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err := s.db.Exec("INSERT INTO post_reactions (post_id, user_id, like_or_dislike) VALUES ($1, $2, $3)", reaction.PostID, reaction.UserID, reaction.LikeOrDislike)
			if err != nil {
				return err
			}

			err = s.UpdatePostReactions(reaction.PostID)
			if err != nil {
				return err
			}

			return nil
		}
		return err
	}

	_, err = s.db.Exec("UPDATE post_reactions SET like_or_dislike = $1 WHERE post_id = $2 AND user_id = $3", reaction.LikeOrDislike, reaction.PostID, reaction.UserID)
	if err != nil {
		return err
	}

	err = s.UpdatePostReactions(reaction.PostID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeletePostReaction(reaction PostReaction) error {
	_, err := s.db.Exec("DELETE FROM post_reactions WHERE post_id = $1 AND user_id = $2", reaction.PostID, reaction.UserID)
	if err != nil {
		return err
	}

	err = s.UpdatePostReactions(reaction.PostID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetPostReaction(post_id, user_id int) (PostReaction, error) {
	var reaction PostReaction
	query := "SELECT post_id, user_id, like_or_dislike FROM post_reactions WHERE post_id = $1 AND user_id = $2"
	err := s.db.QueryRow(query, post_id, user_id).Scan(&reaction.PostID, &reaction.UserID, &reaction.LikeOrDislike)

	return reaction, err
}
