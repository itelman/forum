package handler

import (
	"forum/internal/repository/models"
	"forum/internal/service/auth"
	"forum/internal/service/tmpldata"
	"forum/pkg/forms"
	"net/http"
	"strconv"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		h.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	s, err := h.App.Repository.Posts.Latest()
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	c, err := h.App.Repository.Categories.Latest()
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	err = h.App.Render(w, r, &tmpldata.TemplateData{
		TemplateName: "home_page.html",
		Form:         forms.New(nil),
		Posts:        s,
		Categories:   c,
	})
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}
}

func (h *Handlers) Results(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/results" {
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
	form.RequiredAtLeastOne("categories", "created", "liked")
	if !form.Valid() {
		sesStore.PutSessionData(sessionID, "flash", "Please select at least one filter.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	for _, idStr := range r.Form["categories"] {
		id, err := strconv.Atoi(idStr)
		if err != nil || idStr != strconv.Itoa(id) {
			h.ClientErrorHandler(w, r, http.StatusBadRequest)
			return
		}

		_, err = h.App.Repository.Categories.Get(id)
		if err != nil {
			if err == models.ErrNoRecord {
				h.NotFoundHandler(w, r)
			} else {
				h.ServerErrorHandler(w, r, err)
			}
			return
		}
	}

	created := form.Get("created")
	liked := form.Get("liked")

	if !(created == "" || created == "0" || created == "1") {
		h.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	if !(liked == "" || liked == "0" || liked == "1") {
		h.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	loggedUser := auth.AuthenticatedUser(r)
	if loggedUser == nil && form.ProvidedAtLeastOne("created", "liked") {
		sesStore.PutSessionData(sessionID, "flash", "Your session has expired. Please sign in again.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user_id := -1
	if loggedUser != nil {
		user_id = loggedUser.ID
	}

	s, err := h.App.Repository.Posts.Filter(user_id, r.Form, h.App.Repository.Post_Reactions.FilterByLiked, h.App.Repository.Post_Category.FilterByCategories)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	c, err := h.App.Repository.Categories.Latest()
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	err = h.App.Render(w, r, &tmpldata.TemplateData{
		TemplateName: "home_page.html",
		Form:         forms.New(nil),
		Posts:        s,
		Categories:   c,
	})
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}
}
