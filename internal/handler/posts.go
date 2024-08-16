package handler

import (
	"fmt"
	"forum/internal/config"
	"forum/internal/repository/models"
	"forum/internal/service/forms"
	"net/http"
	"strconv"
)

func (h *Handlers) CreatePostForm(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/create" {
		h.App.NotFoundHandler(w, r)
		return
	}

	if r.Method == http.MethodPost {
		h.CreatePost(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		h.App.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	c, err := h.App.Categories.Latest()
	if err != nil {
		h.App.ServerErrorHandler(w, r, err)
		return
	}

	h.App.Render(w, r, "create_page.html", &config.TemplateData{
		Form:       forms.New(nil),
		Categories: c,
	})
}

func (h *Handlers) CreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.App.ClientErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("user_id", "title", "content", "categories")
	form.MaxLength("title", 100)

	c, err := h.App.Categories.Latest()
	if err != nil {
		h.App.ServerErrorHandler(w, r, err)
		return
	}

	if !form.Valid() {
		h.App.Render(w, r, "create_page.html", &config.TemplateData{
			Form:       form,
			Categories: c,
		})
		return
	}
	id, err := h.App.Posts.Insert(form.Get("user_id"), form.Get("title"), form.Get("content"))
	if err != nil {
		h.App.ServerErrorHandler(w, r, err)
		return
	}

	err = h.App.Post_category.Insert(strconv.Itoa(id), r.PostForm["categories"])
	if err != nil {
		h.App.ServerErrorHandler(w, r, err)
		return
	}
	sessionID, err := h.App.GetSessionIDFromRequest(w, r)
	if err != nil {
		h.App.ServerErrorHandler(w, r, err)
		return
	}
	h.App.PutSessionData(sessionID, "flash", "Post successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", id), http.StatusSeeOther)
}

func (h *Handlers) ShowPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post" {
		h.App.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		h.App.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	idQuery := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idQuery)
	if err != nil || idQuery != strconv.Itoa(id) {
		h.App.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	s, err := h.App.Posts.Get(id)
	if err == models.ErrNoRecord {
		h.App.NotFoundHandler(w, r)
		return
	} else if err != nil {
		h.App.ServerErrorHandler(w, r, err)
		return
	}

	c, err := h.App.Comments.Latest(id)
	if err != nil {
		h.App.ServerErrorHandler(w, r, err)
		return
	}

	categories, err := h.App.Post_category.Get(id)
	if err != nil {
		h.App.ServerErrorHandler(w, r, err)
		return
	}

	loggedUser := h.App.AuthenticatedUser(r)
	if loggedUser != nil {
		reacted, err := h.App.Post_reactions.Get(strconv.Itoa(id), strconv.Itoa(loggedUser.ID))
		if err != nil {
			h.App.ServerErrorHandler(w, r, err)
			return
		}

		s.ReactedByUser = reacted

		for _, comment := range c {
			reacted, err := h.App.Comment_reactions.Get(strconv.Itoa(comment.ID), strconv.Itoa(loggedUser.ID))
			if err != nil {
				h.App.ServerErrorHandler(w, r, err)
				return
			}

			comment.ReactedByUser = reacted
		}
	}

	h.App.Render(w, r, "show_page.html", &config.TemplateData{
		Post:        s,
		Comments:    c,
		PCRelations: categories,
	})
}
