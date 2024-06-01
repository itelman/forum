package models

import "database/sql"

type CommentReaction struct {
	CommentID     int
	UserID        int
	LikeOrDislike int
}

func (s *Storage) InsertCommentReaction(reaction CommentReaction) error {
	if reaction.LikeOrDislike == -1 || reaction.LikeOrDislike == -2 {
		err := s.DeleteCommentReaction(reaction)
		return err
	}

	_, err := s.GetCommentReaction(reaction.CommentID, reaction.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err := s.db.Exec("INSERT INTO comment_reactions (comment_id, user_id, like_or_dislike) VALUES ($1, $2, $3)", reaction.CommentID, reaction.UserID, reaction.LikeOrDislike)
			if err != nil {
				return err
			}

			err = s.UpdateCommentReactions(reaction.CommentID)
			if err != nil {
				return err
			}

			return nil
		}
		return err
	}

	_, err = s.db.Exec("UPDATE comment_reactions SET like_or_dislike = $1 WHERE comment_id = $2 AND user_id = $3", reaction.LikeOrDislike, reaction.CommentID, reaction.UserID)
	if err != nil {
		return err
	}

	err = s.UpdateCommentReactions(reaction.CommentID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteCommentReaction(reaction CommentReaction) error {
	_, err := s.db.Exec("DELETE FROM comment_reactions WHERE comment_id = $1 AND user_id = $2", reaction.CommentID, reaction.UserID)
	if err != nil {
		return err
	}

	err = s.UpdateCommentReactions(reaction.CommentID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetCommentReaction(comment_id, user_id int) (CommentReaction, error) {
	var reaction CommentReaction
	query := "SELECT comment_id, user_id, like_or_dislike FROM comment_reactions WHERE comment_id = $1 AND user_id = $2"
	err := s.db.QueryRow(query, comment_id, user_id).Scan(&reaction.CommentID, &reaction.UserID, &reaction.LikeOrDislike)

	return reaction, err
}
