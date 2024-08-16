package handler

import (
	"forum/internal/config"
	"forum/internal/service/forms"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.App.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		h.App.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	s, err := h.App.Posts.Latest()
	if err != nil {
		h.App.ServerErrorHandler(w, r, err)
		return
	}

	c, err := h.App.Categories.Latest()
	if err != nil {
		h.App.ServerErrorHandler(w, r, err)
		return
	}

	h.App.Render(w, r, "home_page.html", &config.TemplateData{
		Form:       forms.New(nil),
		Posts:      s,
		Categories: c,
	})
}

func (h *Handlers) Results(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/results" {
		h.App.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		h.App.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		h.App.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}
	form := forms.New(r.Form)
	form.RequiredAtLeastOne("categories", "created", "liked")

	if !form.Valid() {
		sessionID, err := h.App.GetSessionIDFromRequest(w, r)
		if err != nil {
			h.App.ServerErrorHandler(w, r, err)
			return
		}
		h.App.PutSessionData(sessionID, "flash", "Please select at least one filter.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	loggedUser := h.App.AuthenticatedUser(r)
	if loggedUser == nil && form.ProvidedAtLeastOne("created", "liked") {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user_id := -1
	if loggedUser != nil {
		user_id = loggedUser.ID
	}

	s, err := h.App.Posts.Filter(user_id, r.Form, h.App.Post_reactions.FilterByLiked, h.App.Post_category.FilterByCategories)
	if err != nil {
		h.App.ServerErrorHandler(w, r, err)
		return
	}

	c, err := h.App.Categories.Latest()
	if err != nil {
		h.App.ServerError(w, r, err)
		return
	}

	h.App.Render(w, r, "home_page.html", &config.TemplateData{
		Posts:      s,
		Categories: c,
	})
}
