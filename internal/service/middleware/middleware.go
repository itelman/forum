package middleware

import (
	"context"
	"fmt"
	"forum/internal/handler"
	"forum/internal/repository/models"
	"forum/internal/service/auth"
	"net/http"
	"sync"
	"time"
)

type Middleware struct {
	Handlers *handler.Handlers

	Limiters        map[string]chan time.Time
	BlockedSessions map[string]time.Time
	Mu              sync.Mutex
}

func (m *Middleware) GetRateLimiter(sessionID string) chan time.Time {
	m.Mu.Lock()
	defer m.Mu.Unlock()

	limiter, exists := m.Limiters[sessionID]
	if !exists {
		// Create a new rate limiter for this session
		limiter = make(chan time.Time, 3)
		for i := 0; i < 3; i++ {
			limiter <- time.Now()
		}

		go func() {
			filler := time.NewTicker(100 * time.Millisecond)
			for t := range filler.C {
				select {
				case limiter <- t:
				default: // If the channel is full, discard the tick
				}
			}
		}()

		m.Limiters[sessionID] = limiter
	}

	return limiter
}

func (m *Middleware) RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sesStore := m.Handlers.App.SessionStore

		sessionID, err := sesStore.GetSessionIDFromRequest(w, r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		m.Mu.Lock()
		lastRequest, exists := m.BlockedSessions[sessionID]
		m.Mu.Unlock()

		if exists {
			if !time.Now().Before(lastRequest.Add(time.Hour)) {
				m.Mu.Lock()
				delete(m.BlockedSessions, sessionID)
				m.Mu.Unlock()
			} else {
				m.Handlers.ClientErrorHandler(w, r, http.StatusTooManyRequests)
				return
			}
		}

		limiter := m.GetRateLimiter(sessionID)

		select {
		case <-limiter:
			next.ServeHTTP(w, r)
		default:
			sesStore.DeleteSession(sessionID)

			m.Mu.Lock()
			m.BlockedSessions[sessionID] = time.Now()
			m.Mu.Unlock()

			m.Handlers.ClientErrorHandler(w, r, http.StatusTooManyRequests)
		}
	})
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

		sesStore := m.Handlers.App.SessionStore
		sessionID, err := sesStore.GetSessionIDFromRequest(w, r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		m.Handlers.App.InfoLog.Printf("COOKIE: %#v\n", sessionID)
		m.Handlers.App.InfoLog.Printf("SESSIONS: %v\n", sesStore.Store)

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RequireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth.AuthenticatedUser(r) == nil {
			http.Redirect(w, r, "/user/login", http.StatusMovedPermanently)
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
		if session == nil || !session.Active {
			next.ServeHTTP(w, r)
			return
		}

		userID := session.UserID
		if userID <= 0 {
			next.ServeHTTP(w, r)
			return
		}

		user, err := m.Handlers.App.Repository.Users.Get(userID)
		if err == models.ErrNoRecord {
			sesStore.DeleteSession(sessionID)
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			m.Handlers.ServerErrorHandler(w, r, err)
			return
		}

		lastRequest := session.LastRequest
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

		/*user, err := m.Handlers.App.Repository.Users.Get(1)
		if err != nil {
			m.Handlers.ServerErrorHandler(w, r, err)
			return
		}*/

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
		return m.RateLimiter(m.Authenticate(m.RequireAuthenticatedUser(nextHandler)))
	}

	return m.RateLimiter(m.Authenticate(nextHandler))
}
