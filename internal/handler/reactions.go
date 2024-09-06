package handler

import (
	"fmt"
	"forum/internal/repository/models"
	"forum/internal/service/auth"
	"forum/pkg/forms"
	"net/http"
	"strconv"
)

func (h *Handlers) HandlePostReaction(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/reaction" {
		h.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		h.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	loggedUser := auth.AuthenticatedUser(r)

	err := r.ParseForm()
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("post_id", "is_like")

	if !form.Valid() {
		h.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	postIdQuery := form.Get("post_id")
	post_id, err := strconv.Atoi(postIdQuery)
	if err != nil || postIdQuery != strconv.Itoa(post_id) {
		h.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	likeQuery := form.Get("is_like")
	is_like, err := strconv.Atoi(likeQuery)
	if err != nil || likeQuery != strconv.Itoa(is_like) {
		h.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	if !(is_like == 0 || is_like == 1) {
		h.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	err = h.App.Repository.Post_Reactions.Insert(post_id, loggedUser.ID, is_like)
	if err != nil {
		if err == models.ErrNoRecord {
			h.NotFoundHandler(w, r)
		} else {
			h.ServerErrorHandler(w, r, err)
		}
		return
	}

	err = h.App.Repository.Posts.UpdateReactions(post_id, h.App.Repository.Post_Reactions.Likes, h.App.Repository.Post_Reactions.Dislikes)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", post_id), http.StatusMovedPermanently)
}

func (h *Handlers) HandleCommentReaction(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comment/reaction" {
		h.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		h.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	loggedUser := auth.AuthenticatedUser(r)

	err := r.ParseForm()
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("comment_id", "is_like")

	if !form.Valid() {
		h.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	commentIdQuery := form.Get("comment_id")

	comment_id, err := strconv.Atoi(commentIdQuery)
	if err != nil || commentIdQuery != strconv.Itoa(comment_id) {
		h.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	likeQuery := form.Get("is_like")
	is_like, err := strconv.Atoi(likeQuery)
	if err != nil || likeQuery != strconv.Itoa(is_like) {
		h.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	if !(is_like == 0 || is_like == 1) {
		h.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	err = h.App.Repository.Comment_Reactions.Insert(comment_id, loggedUser.ID, is_like)
	if err != nil {
		if err == models.ErrNoRecord {
			h.NotFoundHandler(w, r)
		} else {
			h.ServerErrorHandler(w, r, err)
		}
		return
	}

	err = h.App.Repository.Comments.UpdateReactions(comment_id, h.App.Repository.Comment_Reactions.Likes, h.App.Repository.Comment_Reactions.Dislikes)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	comment, err := h.App.Repository.Comments.Get(comment_id)
	if err != nil {
		if err == models.ErrNoRecord {
			h.NotFoundHandler(w, r)
		} else {
			h.ServerErrorHandler(w, r, err)
		}
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", comment.PostID), http.StatusMovedPermanently)
}
