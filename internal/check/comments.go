package check

import (
	"fmt"
	"forum/pkg/forms"
	"net/http"
)

func (app *application) createComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/sneep/create/comment/" {
		app.notFound(w, r)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		app.clientError(w, r, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("snippet_id", "user_id", "content")

	if !form.Valid() {
		app.session.Put(r, "flash", "Please type something into the comment section.")
		http.Redirect(w, r, fmt.Sprintf("/sneep?id=%s", form.Get("snippet_id")), http.StatusSeeOther)
		return
	}

	err = app.comments.Insert(form.Get("snippet_id"), form.Get("user_id"), form.Get("content"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.session.Put(r, "flash", "Comment successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/sneep?id=%s", form.Get("snippet_id")), http.StatusSeeOther)
}
