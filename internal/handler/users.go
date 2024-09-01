package handler

import (
	"fmt"
	"forum/internal/repository/models"
	"forum/internal/service/auth"
	"forum/internal/service/tmpldata"
	"forum/pkg/forms"
	"net/http"
)

func (h *Handlers) SignupUserForm(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/signup" {
		h.NotFoundHandler(w, r)
		return
	}

	if r.Method == http.MethodPost {
		h.SignupUser(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		h.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	sesStore := h.App.SessionStore

	if auth.AuthenticatedUser(r) != nil {
		sessionID, err := sesStore.GetSessionIDFromRequest(w, r)
		if err != nil {
			h.ServerErrorHandler(w, r, err)
			return
		}

		sesStore.PutSessionData(sessionID, "flash", "Please log out before proceeding.")

		err = h.App.Render(w, r, &tmpldata.TemplateData{
			TemplateName: "signup_page.html",
		})
		if err != nil {
			h.ServerErrorHandler(w, r, err)
			return
		}

		return
	}

	err := h.App.Render(w, r, &tmpldata.TemplateData{
		TemplateName: "signup_page.html",
		Form:         forms.New(nil),
	})
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}
}

func (h *Handlers) SignupUser(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/signup" {
		h.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		h.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")

	form.MinLength("name", 5)
	form.MaxLength("name", 15)

	form.MinLength("email", 15)
	form.MaxLength("email", 20)

	form.MinLength("password", 6)
	form.MaxLength("password", 10)

	form.MatchesPattern("name", forms.NameRX)
	form.MatchesPattern("email", forms.EmailRX)
	form.MatchesPattern("password", forms.PasswordRX)

	if !form.Valid() {
		err = h.App.Render(w, r, &tmpldata.TemplateData{
			TemplateName: "signup_page.html",
			Form:         form,
		})
		if err != nil {
			h.ServerErrorHandler(w, r, err)
			return
		}

		return
	}
	err = h.App.Repository.Users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err == models.ErrDuplicateNameOrEmail {
		form.Errors.Add("generic", "Username and/or email are already in use")

		err = h.App.Render(w, r, &tmpldata.TemplateData{
			TemplateName: "signup_page.html",
			Form:         form,
		})
		if err != nil {
			h.ServerErrorHandler(w, r, err)
			return
		}

		return
	} else if err != nil {
		fmt.Println(err)
		fmt.Printf("**%s**", err)
		h.ServerErrorHandler(w, r, err)
		return
	}

	sesStore := h.App.SessionStore

	sessionID, err := sesStore.GetSessionIDFromRequest(w, r)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}
	sesStore.PutSessionData(sessionID, "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (h *Handlers) LoginUserForm(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/login" {
		h.NotFoundHandler(w, r)
		return
	}
	if r.Method == http.MethodPost {
		h.LoginUser(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		h.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	sesStore := h.App.SessionStore

	if auth.AuthenticatedUser(r) != nil {
		sessionID, err := sesStore.GetSessionIDFromRequest(w, r)
		if err != nil {
			h.ServerErrorHandler(w, r, err)
			return
		}

		sesStore.PutSessionData(sessionID, "flash", "Please log out before proceeding.")

		err = h.App.Render(w, r, &tmpldata.TemplateData{
			TemplateName: "login_page.html",
		})
		if err != nil {
			h.ServerErrorHandler(w, r, err)
			return
		}

		return
	}

	err := h.App.Render(w, r, &tmpldata.TemplateData{
		TemplateName: "login_page.html",
		Form:         forms.New(nil),
	})
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}
}

func (h *Handlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/login" {
		h.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		h.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}
	form := forms.New(r.PostForm)
	id, err := h.App.Repository.Users.Authenticate(form.Get("name"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Username and/or Password are incorrect")

		err = h.App.Render(w, r, &tmpldata.TemplateData{
			TemplateName: "login_page.html",
			Form:         form,
		})
		if err != nil {
			h.ServerErrorHandler(w, r, err)
			return
		}

		return
	} else if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	sesStore := h.App.SessionStore

	existingSessionID, exists := sesStore.GetSessionByUserID(id)
	if exists {
		sesStore.DisableSession(existingSessionID)
		sesStore.PutSessionData(existingSessionID, "flash", "Your session has expired. Please sign in again.")
	}

	sessionID, err := sesStore.CreateNewSession(id)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}
	sesStore.PutSessionData(sessionID, "userID", id)

	http.SetCookie(w, &http.Cookie{
		Name:  "session_id",
		Value: sessionID,
		Path:  "/",
		// MaxAge: int(time.Duration(h.App.CookieLimit).Seconds()),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handlers) LogoutUser(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/logout" {
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

	sesStore.DisableSession(sessionID)
	sesStore.PutSessionData(sessionID, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
