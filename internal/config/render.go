package config

import (
	"bytes"
	"fmt"
	"forum/internal/service/auth"
	"forum/internal/service/tmpldata"
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
	session := sesStore.GetSession(sessionID)

	td.AddDefaultData(auth.AuthenticatedUser(r), sesStore.PopSessionFlash(sessionID))

	if !session.Active {
		sesStore.DeleteSession(sessionID)
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
