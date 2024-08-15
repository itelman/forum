package config

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"forum/internal/repository/models"
	"net/http"
	"runtime/debug"
	"time"
)

func GenerateCSRFToken() (string, error) {
	tokenBytes := make([]byte, 32) // 32 bytes = 256 bits
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(tokenBytes), nil
}

func (app *Application) AddDefaultData(w http.ResponseWriter, td *TemplateData, r *http.Request) *TemplateData {
	if td == nil {
		td = &TemplateData{}
	}
	sessionID, err := app.GetSessionIDFromRequest(w, r)
	if err == nil {
		csrfToken, _ := app.GetSessionToken(sessionID)
		td.CSRFToken = csrfToken
	} else {
		td.CSRFToken = ""
	}

	td.AuthenticatedUser = app.AuthenticatedUser(r)
	td.CurrentYear = time.Now().Year()

	if err == nil {
		flash := app.GetSession(sessionID)["flash"]
		if flash != nil {
			td.Flash = flash.(string)
			app.DeleteSessionData(sessionID, "flash")
		} else {
			td.Flash = ""
		}
	}

	return td
}

func (app *Application) Render(w http.ResponseWriter, r *http.Request, name string, td *TemplateData) {
	ts, ok := app.Templatecache[name]
	if !ok {
		if name == "error_page.html" {
			app.ServerError(w, r, fmt.Errorf("The template %s does not exist", name))
		} else {
			app.ServerErrorHandler(w, r, fmt.Errorf("The template %s does not exist", name))
		}
		return
	}

	buf := new(bytes.Buffer)
	err := ts.Execute(buf, app.AddDefaultData(w, td, r))
	if err != nil {
		if name == "error_page.html" {
			app.ServerError(w, r, err)
		} else {
			app.ServerErrorHandler(w, r, err)
		}
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		panic(err)
	}
}

func (app *Application) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) AuthenticatedUser(r *http.Request) *models.User {
	value := r.Context().Value(contextKeyUser)

	user, ok := value.(*models.User)
	if !ok {
		return nil
	}

	return user
}
