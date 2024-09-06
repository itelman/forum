package handler

import (
	"forum/internal/repository/models"
	"forum/internal/service/auth"
	"forum/internal/service/tmpldata"
	"forum/pkg/forms"
	"net/http"
)

func (h *Handlers) ShowNotifications(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/notifications" {
		h.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		h.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	sesStore := h.App.SessionStore

	sessionID, err := sesStore.GetSessionIDFromRequest(w, r)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	form := forms.New(r.Form)
	form.Required("filter")
	form.PermittedValues("filter", "1", "2")
	if !form.Valid() {
		sesStore.PutSessionData(sessionID, "flash", "Please select one filter.")
		http.Redirect(w, r, "/user/notifications?filter=1", http.StatusMovedPermanently)
		return
	}

	if form.Get("filter") == "1" {
		h.ShowComments(w, r)
		return
	} else {
		h.ShowReactions(w, r)
		return
	}
}

func (h *Handlers) ShowComments(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/notifications" {
		h.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		h.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	loggedUser := auth.AuthenticatedUser(r)

	posts, err := h.App.Repository.Posts.Created(loggedUser.ID)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	var commentsAll []*models.Comment

	for _, post_comments := range posts {
		comments, err := h.App.Repository.Comments.LatestIgnoreUser(post_comments.Post.ID, loggedUser.ID)
		if err != nil {
			h.ServerErrorHandler(w, r, err)
			return
		}

		commentsAll = append(commentsAll, comments...)
	}

	err = h.App.Render(w, r, &tmpldata.TemplateData{
		TemplateName: "notifications_page.html",
		Comments:     commentsAll,
	})
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}
}

func (h *Handlers) ShowReactions(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/notifications" {
		h.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		h.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	loggedUser := auth.AuthenticatedUser(r)

	posts, err := h.App.Repository.Posts.Created(loggedUser.ID)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	var reactionsAll []*models.PostReaction

	for _, post_comments := range posts {
		post_reactions, err := h.App.Repository.Post_Reactions.LatestIgnoreUser(post_comments.Post.ID, loggedUser.ID)
		if err != nil {
			h.ServerErrorHandler(w, r, err)
			return
		}

		reactionsAll = append(reactionsAll, post_reactions...)
	}

	err = h.App.Render(w, r, &tmpldata.TemplateData{
		TemplateName:   "notifications_page.html",
		Post_Reactions: reactionsAll,
	})
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}
}
