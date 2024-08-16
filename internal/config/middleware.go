package config

import (
	"context"
	"fmt"
	"forum/internal/repository/models"
	"net/http"
	"time"
)

func (app *Application) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.ServerErrorHandler(w, r, fmt.Errorf("%s", err))
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

func (app *Application) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *Application) RequireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*if app.AuthenticatedUser(r) == nil {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}*/

		next.ServeHTTP(w, r)
	})
}

func (app *Application) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := app.GetSessionIDFromRequest(w, r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		sessionData := app.GetSession(sessionID)
		/*fmt.Println("COOKIE: ", sessionData)
		fmt.Println("SESSIONS: ", app.SessionStore)*/

		_, exists := sessionData["inactive"]
		if exists {
			app.DeleteSession(sessionID)
			app.PutSessionData(sessionID, "flash", "You've been logged out. Please sign in again.")
			next.ServeHTTP(w, r)
			return
		}

		flash, exists := sessionData["flash"]
		if exists && flash == "You've been logged out successfully!" {
			app.PutSessionData(sessionID, "inactive", 1)
			next.ServeHTTP(w, r)
			return
		}

		userID, exists := sessionData["userID"]
		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		user, err := app.Users.Get(userID.(int))
		if err == models.ErrNoRecord {
			app.DeleteSession(sessionID)
			delete(app.ActiveSessions, userID.(int))
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.ServerErrorHandler(w, r, err)
			return
		}

		lastRequest, exists := sessionData["lastRequest"].(time.Time)
		if !exists {
			app.ServerErrorHandler(w, r, err)
			return
		}

		if !time.Now().Before(lastRequest.Add(app.CookieLimit)) {
			app.DeleteSession(sessionID)
			delete(app.ActiveSessions, sessionData["userID"].(int))

			app.PutSessionData(sessionID, "flash", "Your session has expired. Please sign in again.")
			next.ServeHTTP(w, r)
			return
		}

		app.UpdateSessionLastReq(sessionID)

		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *Application) StandardMiddleware(next http.Handler) http.Handler {
	return app.RecoverPanic(app.LogRequest(SecureHeaders(next)))
}

func (app *Application) DynamicMiddleware(next func(http.ResponseWriter, *http.Request), requireAuth bool) http.Handler {
	nextHandler := http.HandlerFunc(next)

	if requireAuth {
		return app.Authenticate(app.RequireAuthenticatedUser(nextHandler))
	}

	return app.Authenticate(nextHandler)
}
