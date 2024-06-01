package handlers

import (
	"forum/internal/models"
	"net/http"
	"strconv"
)

func CommentVerifyHandler(w http.ResponseWriter, r *http.Request, db *models.Storage, user_id int) {
	handlerName := "CommentVerifyHandler"

	comment := models.Comment{}

	for k, v := range r.Form {
		switch k {
		case "id":
			post_id, err := strconv.Atoi(v[0])
			if err != nil {
				ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
				return
			}
			comment.PostID = post_id
		case "comment":
			comment.Content = v[0]
		}
	}
	comment.UserID = user_id

	err := db.InsertComment(comment)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}
}
