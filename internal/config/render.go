package config

import (
	"bytes"
	"fmt"
	"forum/internal/service/auth"
	"forum/internal/service/tmpldata"
	"forum/pkg/csrf"
	"net/http"
)

func (app *Application) Render(w http.ResponseWriter, r *http.Request, td *tmpldata.TemplateData) error {
	ts, ok := app.TemplateCache[td.TemplateName]
	if !ok {
		return fmt.Errorf("the template %s does not exist", td.TemplateName)
	}

	sesStore := app.SessionStore

	sessionID, err := sesStore.GetSessionIDFromRequest(w, r)
	if err != nil {
		return err
	}

	td.AddDefaultData(auth.AuthenticatedUser(r), sesStore.PopSessionFlash(sessionID))

	csrfToken, err := csrf.NewToken()
	if err != nil {
		return err
	}
	td.CSRFToken = csrfToken
	sesStore.PutSessionData(sessionID, "csrf_token", td.CSRFToken)

	session := sesStore.GetSession(sessionID)
	if session != nil && !session.Active {
		sesStore.DeleteSession(sessionID)

		newSesID, err := sesStore.NewSessionID()
		if err != nil {
			return err
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    newSesID,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			// MaxAge: int(time.Duration(h.App.CookieLimit).Seconds()),
		})
	}

	buf := new(bytes.Buffer)
	err = ts.Execute(buf, td)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}
