package check

import (
	"fmt"
	"forum/pkg/forms"
	"forum/pkg/models"
	"net/http"
	"strconv"
)

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/sneep/create" {
		app.notFound(w, r)
		return
	}

	if r.Method == http.MethodPost {
		app.createSnippet(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		app.clientError(w, r, http.StatusMethodNotAllowed)
		return
	}

	c, err := app.categories.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(w, r, "create_page.html", &templateData{
		Form:       forms.New(nil),
		Categories: c,
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("user_id", "title", "content", "categories")
	form.MaxLength("title", 100)

	c, err := app.categories.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if !form.Valid() {
		app.render(w, r, "create_page.html", &templateData{
			Form:       form,
			Categories: c,
		})
		return
	}
	id, err := app.snippets.Insert(form.Get("user_id"), form.Get("title"), form.Get("content"), form.Get("likes"), form.Get("dislikes"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.post_category.Insert(strconv.Itoa(id), r.PostForm["categories"])
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.session.Put(r, "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/sneep?id=%d", id), http.StatusSeeOther)
}

func (app *application) snippet(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/sneep" {
		app.notFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		app.clientError(w, r, http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w, r)
		return
	}
	s, err := app.snippets.Get(int(id))
	if err == models.ErrNoRecord {
		app.notFound(w, r)
		return
	} else if err != nil {
		app.serverError(w, r, err)
		return
	}

	c, err := app.comments.Latest(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	categories, err := app.post_category.Get(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(w, r, "show_page.html", &templateData{
		Snippet:     s,
		Comments:    c,
		PCRelations: categories,
	})
}
