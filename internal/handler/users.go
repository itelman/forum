package handler

import (
	"fmt"
	"forum/internal/config"
	"forum/internal/repository/models"
	"forum/internal/service/forms"
	"net/http"
)

func (h *Handlers) SignupUserForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.SignupUser(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		h.App.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	h.App.Render(w, r, "signup_page.html", &config.TemplateData{
		Form: forms.New(nil),
	})
}

func (h *Handlers) SignupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.App.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 6)

	if !form.Valid() {
		h.App.Render(w, r, "signup_page.html", &config.TemplateData{
			Form: form,
		})
		return
	}
	err = h.App.Users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err == models.ErrDuplicateNameOrEmail {
		form.Errors.Add("generic", "Username and/or email are already in use")
		h.App.Render(w, r, "signup_page.html", &config.TemplateData{
			Form: form,
		})
		return
	} else if err != nil {
		fmt.Println(err)
		fmt.Printf("**%s**", err)
		h.App.ServerErrorHandler(w, r, err)
		return
	}

	sessionID, err := h.App.GetSessionIDFromRequest(w, r)
	if err != nil {
		h.App.ServerErrorHandler(w, r, err)
		return
	}
	h.App.PutSessionData(sessionID, "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (h *Handlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.App.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	id, err := h.App.Users.Authenticate(form.Get("name"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Username and/or Password are incorrect")
		h.App.Render(w, r, "login_page.html", &config.TemplateData{
			Form: form,
		})
		return
	} else if err != nil {
		h.App.ServerErrorHandler(w, r, err)
		return
	}
	sessionID, err := h.App.CreateNewSession(id)
	if err != nil {
		h.App.ServerErrorHandler(w, r, err)
		return
	}
	h.App.PutSessionData(sessionID, "userID", id)
	http.SetCookie(w, &http.Cookie{
		Name:  "session_id",
		Value: sessionID,
		Path:  "/",
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handlers) LoginUserForm(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/login" {
		h.App.NotFoundHandler(w, r)
		return
	}
	if r.Method == http.MethodPost {
		h.LoginUser(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		h.App.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}
	h.App.Render(w, r, "login_page.html", &config.TemplateData{
		Form: forms.New(nil),
	})
}

func (h *Handlers) LogoutUser(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/logout" {
		h.App.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		h.App.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}
	sessionID, err := h.App.GetSessionIDFromRequest(w, r)
	if err != nil {
		h.App.ServerErrorHandler(w, r, err)
		return
	}
	userID := h.App.GetSessionUserID(sessionID)
	h.App.DeleteSession(sessionID)
	delete(h.App.ActiveSessions, userID)
	h.App.PutSessionData(sessionID, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
