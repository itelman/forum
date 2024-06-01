package handlers

import (
	"forum/internal/models"
	"net/http"
	"strconv"
)

func PostReactionHandler(w http.ResponseWriter, r *http.Request, db *models.Storage, user_id int) {
	handlerName := "PostReactionHandler"

	reaction := models.PostReaction{}

	for k, v := range r.Form {
		switch k {
		case "id":
			post_id, err := strconv.Atoi(v[0])
			if err != nil {
				ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
				return
			}
			reaction.PostID = post_id
		case "postReaction":
			like_or_dislike, err := strconv.Atoi(v[0])
			if err != nil {
				ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
				return
			}
			reaction.LikeOrDislike = like_or_dislike
		}
	}
	reaction.UserID = user_id

	err := db.InsertPostReaction(reaction)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}
}

func CommentReactionHandler(w http.ResponseWriter, r *http.Request, db *models.Storage, user_id int) {
	handlerName := "CommentReactionHandler"

	reaction := models.CommentReaction{}

	for k, v := range r.Form {
		switch k {
		case "comment_id":
			comment_id, err := strconv.Atoi(v[0])
			if err != nil {
				ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
				return
			}
			reaction.CommentID = comment_id
		case "commentReaction":
			like_or_dislike, err := strconv.Atoi(v[0])
			if err != nil {
				ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
				return
			}
			reaction.LikeOrDislike = like_or_dislike
		}
	}
	reaction.UserID = user_id

	err := db.InsertCommentReaction(reaction)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}
}
