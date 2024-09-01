package middleware

import (
	"context"
	"fmt"
	"forum/internal/handler"
	"forum/internal/repository/models"
	"forum/internal/service/auth"
	"net/http"
	"time"
)

type Middleware struct {
	Handlers *handler.Handlers
}

func (m *Middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				m.Handlers.ServerErrorHandler(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.Handlers.App.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RequireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth.AuthenticatedUser(r) == nil {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sesStore := m.Handlers.App.SessionStore

		sessionID, err := sesStore.GetSessionIDFromRequest(w, r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		session := sesStore.GetSession(sessionID)
		if !session.Active {
			next.ServeHTTP(w, r)
			return
		}

		userID := sesStore.GetSessionUserID(sessionID)
		if userID <= 0 {
			next.ServeHTTP(w, r)
			return
		}

		/*fmt.Printf("COOKIE: %#v\n", sessionID)
		fmt.Printf("SESSIONS: %#v\n", sesStore.Store)*/

		user, err := m.Handlers.App.Repository.Users.Get(userID)
		if err == models.ErrNoRecord {
			sesStore.DeleteSession(sessionID)
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			m.Handlers.ServerErrorHandler(w, r, err)
			return
		}

		sessionData := sesStore.GetSession(sessionID)
		lastRequest := sessionData.LastRequest
		if lastRequest.IsZero() {
			m.Handlers.ServerErrorHandler(w, r, err)
			return
		}

		if !time.Now().Before(lastRequest.Add(m.Handlers.App.CookieLimit)) {
			sesStore.DisableSession(sessionID)
			sesStore.PutSessionData(sessionID, "flash", "Your session has expired. Please sign in again.")
			next.ServeHTTP(w, r)
			return
		}

		sesStore.UpdateSessionLastReq(sessionID)

		ctx := context.WithValue(r.Context(), auth.ContextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) StandardMiddleware(next http.Handler) http.Handler {
	return m.RecoverPanic(m.LogRequest(SecureHeaders(next)))
}

func (m *Middleware) DynamicMiddleware(next func(http.ResponseWriter, *http.Request), requireAuth bool) http.Handler {
	nextHandler := http.HandlerFunc(next)

	if requireAuth {
		return m.Authenticate(m.RequireAuthenticatedUser(nextHandler))
	}

	return m.Authenticate(nextHandler)
}
