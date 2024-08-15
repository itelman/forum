package handler

import (
	"fmt"
	"forum/internal/service/forms"
	"net/http"
	"strconv"
)

func (h *Handlers) HandlePostReaction(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/reaction" {
		h.App.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		h.App.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		h.App.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("post_id", "user_id", "is_like")

	if !form.Valid() {
		h.App.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	post_id, err := strconv.Atoi(form.Get("post_id"))
	if err != nil || post_id < 1 {
		h.App.NotFoundHandler(w, r)
		return
	}

	err = h.App.Post_reactions.Insert(form.Get("post_id"), form.Get("user_id"), form.Get("is_like"))
	if err != nil {
		h.App.ServerErrorHandler(w, r, err)
		return
	}

	err = h.App.Posts.UpdateReactions(post_id, h.App.Post_reactions.Likes, h.App.Post_reactions.Dislikes)
	if err != nil {
		h.App.ServerErrorHandler(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post?id=%s", form.Get("post_id")), http.StatusSeeOther)
}

func (h *Handlers) HandleCommentReaction(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comment/reaction" {
		h.App.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		h.App.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		h.App.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("post_id", "comment_id", "user_id", "is_like")

	if !form.Valid() {
		h.App.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	comment_id, err := strconv.Atoi(form.Get("comment_id"))
	if err != nil || comment_id < 1 {
		h.App.NotFoundHandler(w, r)
		return
	}

	err = h.App.Comment_reactions.Insert(form.Get("comment_id"), form.Get("user_id"), form.Get("is_like"))
	if err != nil {
		h.App.ServerErrorHandler(w, r, err)
		return
	}

	err = h.App.Comments.UpdateReactions(comment_id, h.App.Comment_reactions.Likes, h.App.Comment_reactions.Dislikes)
	if err != nil {
		h.App.ServerErrorHandler(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post?id=%s", form.Get("post_id")), http.StatusSeeOther)
}
