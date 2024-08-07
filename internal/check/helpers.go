package check

import (
	"bytes"
	"fmt"
	"forum/pkg/models"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"
)

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CSRFToken = nosurf.Token(r)
	td.AuthenticatedUser = app.authenticatedUser(r)
	td.CurrentYear = time.Now().Year()
	td.Flash = app.session.PopString(r, "flash")
	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templatecache[name]
	if !ok {
		app.serverError(w, r, fmt.Errorf("The template %s does not exist", name))
		return
	}
	buf := new(bytes.Buffer)
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	buf.WriteTo(w)
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	error := &models.Error{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
	}

	app.render(w, r, "error_page.html", &templateData{
		Error: error,
	})

	// http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, r *http.Request, status int) {
	err := &models.Error{
		Code:    status,
		Message: http.StatusText(status),
	}

	app.render(w, r, "error_page.html", &templateData{
		Error: err,
	})

	// http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	app.clientError(w, r, http.StatusNotFound)
}

func (app *application) authenticatedUser(r *http.Request) *models.User {
	value := r.Context().Value(contextKeyUser)

	user, ok := value.(*models.User)
	if !ok {
		return nil
	}

	return user
}
