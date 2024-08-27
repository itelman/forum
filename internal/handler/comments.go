package handler

import (
	"fmt"
	"forum/internal/service/forms"
	"net/http"
	"strconv"
)

func (h *Handlers) CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/comment" {
		h.App.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		h.App.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	loggedUser := h.App.AuthenticatedUser(r)

	err := r.ParseForm()
	if err != nil {
		h.App.ServerErrorHandler(w, r, err)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("post_id", "content")
	post_id := form.Get("post_id")

	if !form.Valid() {
		sessionID, err := h.App.GetSessionIDFromRequest(w, r)
		if err != nil {
			h.App.ServerErrorHandler(w, r, err)
			return
		}
		h.App.PutSessionData(sessionID, "flash", "Please type something into the comments section.")
		http.Redirect(w, r, fmt.Sprintf("/post?id=%s", post_id), http.StatusSeeOther)
		return
	}

	err = h.App.Comments.Insert(post_id, strconv.Itoa(loggedUser.ID), form.Get("content"))
	if err != nil {
		h.App.ServerErrorHandler(w, r, err)
		return
	}
	sessionID, err := h.App.GetSessionIDFromRequest(w, r)
	if err != nil {
		h.App.ServerErrorHandler(w, r, err)
		return
	}
	h.App.PutSessionData(sessionID, "flash", "Comment successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/post?id=%s", post_id), http.StatusSeeOther)
}
