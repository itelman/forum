package handler

import (
	"fmt"
	"forum/internal/repository/models"
	"forum/internal/service/auth"
	"forum/pkg/forms"
	"net/http"
	"strconv"
)

func (h *Handlers) CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/comment" {
		h.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		h.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	sesStore := h.App.SessionStore

	sessionID, err := sesStore.GetSessionIDFromRequest(w, r)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	loggedUser := auth.AuthenticatedUser(r)

	err = r.ParseForm()
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("post_id", "content")

	postIdQuery := form.Get("post_id")
	post_id, err := strconv.Atoi(postIdQuery)
	if err != nil || postIdQuery != strconv.Itoa(post_id) {
		h.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	if !form.Valid() {
		sesStore.PutSessionData(sessionID, "flash", "Please type something into the comments section.")
		http.Redirect(w, r, fmt.Sprintf("/post?id=%d", post_id), http.StatusMovedPermanently)
		return
	}

	err = h.App.Repository.Comments.Insert(post_id, loggedUser.ID, form.Get("content"))
	if err != nil {
		if err == models.ErrNoRecord {
			h.NotFoundHandler(w, r)
		} else {
			h.ServerErrorHandler(w, r, err)
		}
		return
	}

	sesStore.PutSessionData(sessionID, "flash", "Comment successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", post_id), http.StatusMovedPermanently)
}
