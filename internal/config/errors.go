package config

import (
	"fmt"
	"forum/internal/repository/models"
	"net/http"
	"runtime/debug"
)

func (app *Application) ServerErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)

	errorModel := &models.Error{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
	}

	w.WriteHeader(http.StatusInternalServerError)
	app.Render(w, r, "error_page.html", &TemplateData{
		Error: errorModel,
	})
}

func (app *Application) ClientErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	err := &models.Error{
		Code:    status,
		Message: http.StatusText(status),
	}

	w.WriteHeader(status)
	app.Render(w, r, "error_page.html", &TemplateData{
		Error: err,
	})
}

func (app *Application) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	app.ClientErrorHandler(w, r, http.StatusNotFound)
}
