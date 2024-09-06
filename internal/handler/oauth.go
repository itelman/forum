package handler

import (
	"fmt"
	"forum/internal/repository/models"
	"forum/internal/service/auth"
	"forum/internal/service/auth/github"
	"forum/internal/service/auth/google"
	"forum/internal/service/tmpldata"
	"forum/pkg/forms"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

func (h *Handlers) LoginGithub(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/github" {
		h.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		h.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	if auth.AuthenticatedUser(r) != nil {
		sessionID, err := h.App.SessionStore.GetSessionIDFromRequest(w, r)
		if err != nil {
			h.ServerErrorHandler(w, r, err)
			return
		}

		h.App.SessionStore.PutSessionData(sessionID, "flash", "Please log out before proceeding.")

		err = h.App.Render(w, r, &tmpldata.TemplateData{
			TemplateName: "login_page.html",
		})
		if err != nil {
			h.ServerErrorHandler(w, r, err)
			return
		}

		return
	}

	githubClientID := github.GetClientID()
	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s", githubClientID, "https://forum-099y.onrender.com/auth/github/callback")

	http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
}

func (h *Handlers) LoginGithubCallback(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/github/callback" {
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

	code := r.URL.Query().Get("code")
	authErr := r.URL.Query().Get("error")

	if code == "" {
		if authErr == "" {
			h.ClientErrorHandler(w, r, http.StatusForbidden)
		} else {
			sesStore.PutSessionData(sessionID, "flash", "Authorization unsuccessful. Please try again.")
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		}
		return
	}

	accessToken, err := github.GetAccessToken(code)
	if err != nil {
		if err == models.ErrBadGateway {
			sesStore.PutSessionData(sessionID, "flash", "Authorization unsuccessful. Please try again.")
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		} else {
			h.ServerErrorHandler(w, r, err)
		}
		return
	}

	userData, err := github.GetUserData(accessToken)
	if err != nil {
		if err == models.ErrBadGateway {
			sesStore.PutSessionData(sessionID, "flash", "Authorization unsuccessful. Please try again.")
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		} else {
			h.ServerErrorHandler(w, r, err)
		}
		return
	}

	id, err := h.App.Repository.Users.OAuth(userData.Email, userData.AccountID, userData.Provider)
	if err != nil {
		if err == models.ErrNoRecord {
			sessionID, err := sesStore.GetSessionIDFromRequest(w, r)
			if err != nil {
				h.ServerErrorHandler(w, r, err)
				return
			}

			userData.Username += "_" + strconv.Itoa(rand.Intn(9999-1000)+1000)
			sesStore.PutSessionData(sessionID, "authData", userData)
			http.Redirect(w, r, "/user/signup/provider", http.StatusMovedPermanently)
		} else {
			h.ServerErrorHandler(w, r, err)
		}
		return
	}

	sesStore.DeleteSession(sessionID)

	existingSessionID, exists := sesStore.GetSessionByUserID(id)
	if exists {
		sesStore.DisableSession(existingSessionID)
		sesStore.PutSessionData(existingSessionID, "flash", "Your session has expired. Please sign in again.")
	}

	sessionID, err = sesStore.CreateNewSession(id)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		// MaxAge: int(time.Duration(sesStore.CookieLimit).Seconds()),
	})

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (h *Handlers) LoginGoogle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/google" {
		h.NotFoundHandler(w, r)
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
}

func (h *Handlers) LoginGoogleCallback(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/google/callback" {
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

	code := r.URL.Query().Get("code")
	authErr := r.URL.Query().Get("error")

	if code == "" {
		if authErr == "" {
			h.ClientErrorHandler(w, r, http.StatusForbidden)
		} else {
			sesStore.PutSessionData(sessionID, "flash", "Authorization unsuccessful. Please try again.")
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		}
		return
	}

	accessToken, tokenType, err := google.GetAccessToken(code, "authorization_code")
	if err != nil {
		if err == models.ErrBadGateway {
			sesStore.PutSessionData(sessionID, "flash", "Authorization unsuccessful. Please try again.")
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		} else {
			h.ServerErrorHandler(w, r, err)
		}
		return
	}

	userData, err := google.GetUserData(tokenType, accessToken)
	if err != nil {
		if err == models.ErrBadGateway {
			sesStore.PutSessionData(sessionID, "flash", "Authorization unsuccessful. Please try again.")
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		} else {
			h.ServerErrorHandler(w, r, err)
		}
		return
	}

	id, err := h.App.Repository.Users.OAuth(userData.Email, userData.AccountID, userData.Provider)
	if err != nil {
		if err == models.ErrNoRecord {
			sessionID, err := sesStore.GetSessionIDFromRequest(w, r)
			if err != nil {
				h.ServerErrorHandler(w, r, err)
				return
			}

			// w.Header().Set("Content-Type", "application/json")
			// json.NewEncoder(w).Encode(userData)

			userData.Username += "_" + strconv.Itoa(rand.Intn(9999-1000)+1000)
			sesStore.PutSessionData(sessionID, "authData", userData)
			http.Redirect(w, r, "/user/signup/provider", http.StatusMovedPermanently)
		} else {
			h.ServerErrorHandler(w, r, err)
		}
		return
	}

	sesStore.DeleteSession(sessionID)

	existingSessionID, exists := sesStore.GetSessionByUserID(id)
	if exists {
		sesStore.DisableSession(existingSessionID)
		sesStore.PutSessionData(existingSessionID, "flash", "Your session has expired. Please sign in again.")
	}

	sessionID, err = sesStore.CreateNewSession(id)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		// MaxAge: int(time.Duration(sesStore.CookieLimit).Seconds()),
	})

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (h *Handlers) SignupUserProviderForm(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/signup/provider" {
		h.NotFoundHandler(w, r)
		return
	}

	if r.Method == http.MethodPost {
		h.SignupUserProvider(w, r)
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

	sessionData := sesStore.GetSession(sessionID)
	userData := sessionData.AuthData
	if userData == nil {
		h.ClientErrorHandler(w, r, http.StatusForbidden)
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(userData)

	postForm := make(url.Values)
	postForm.Set("name", userData.Username)
	postForm.Set("email", userData.Email)

	form := forms.New(postForm)
	form.Required("name")

	form.MinLength("name", 5)
	form.MaxLength("name", 30)

	form.MatchesPattern("name", forms.NameRX)

	err = h.App.Render(w, r, &tmpldata.TemplateData{
		TemplateName: "signup_prov_page.html",
		Form:         form,
	})
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}
}

func (h *Handlers) SignupUserProvider(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/signup/provider" {
		h.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		h.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	sesStore := h.App.SessionStore

	sessionID, err := sesStore.GetSessionIDFromRequest(w, r)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	sessionData := sesStore.GetSession(sessionID)
	userData := sessionData.AuthData
	if userData == nil {
		h.ClientErrorHandler(w, r, http.StatusForbidden)
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(userData)

	err = r.ParseForm()
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name")

	form.MinLength("name", 5)
	form.MaxLength("name", 30)

	form.MatchesPattern("name", forms.NameRX)

	if !form.Valid() {
		err = h.App.Render(w, r, &tmpldata.TemplateData{
			TemplateName: "signup_prov_page.html",
			Form:         form,
		})
		if err != nil {
			h.ServerErrorHandler(w, r, err)
			return
		}
		return
	}

	id, err := h.App.Repository.Users.InsertOAuth(form.Get("name"), userData.Email, userData.AccountID, userData.Provider)
	if err == models.ErrDuplicateNameOrEmail {
		form.Errors.Add("generic", "Username and/or email are already in use")

		err = h.App.Render(w, r, &tmpldata.TemplateData{
			TemplateName: "signup_prov_page.html",
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

	sesStore.DeleteSession(sessionID)

	existingSessionID, exists := sesStore.GetSessionByUserID(id)
	if exists {
		sesStore.DisableSession(existingSessionID)
		sesStore.PutSessionData(existingSessionID, "flash", "Your session has expired. Please sign in again.")
	}

	sesStore.DeleteSession(sessionID)

	sessionID, err = sesStore.CreateNewSession(id)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		// MaxAge: int(time.Duration(sesStore.CookieLimit).Seconds()),
	})

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
