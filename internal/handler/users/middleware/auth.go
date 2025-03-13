package middleware

import (
	"context"
	"errors"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/exception"
	"github.com/itelman/forum/internal/service/users"
	"github.com/itelman/forum/pkg/sesm"
	"net/http"
	"sync"
	"time"
)

var (
	sesTimeLimit = 24 * time.Hour
	sesBlockTime = time.Hour
	maxRequests  = 10
)

var (
	ErrTooManyRequests = errors.New("too many requests")
	ErrSessionExpired  = errors.New("session expired")
)

type AuthMiddleware interface {
	Authenticate(next http.Handler) http.Handler
}

type middleware struct {
	users      users.Service
	sesManager sesm.SessionManager
	exceptions exception.Exceptions

	limiters map[string]chan time.Time
	mutex    sync.RWMutex
}

func NewMiddleware(users users.Service, sesManager sesm.SessionManager, exceptions exception.Exceptions) *middleware {
	return &middleware{
		users:      users,
		sesManager: sesManager,
		exceptions: exceptions,
		limiters:   make(map[string]chan time.Time),
	}
}

func (m *middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errInternalSrvResp := func(err error) {
			m.sesManager.DeleteCurrentSession(r)
			http.SetCookie(w, dto.DeleteCookie(sesm.SessionId))
			m.exceptions.ErrInternalServerHandler(w, r, err)
		}

		userId, err := m.sesManager.GetSessionData(r, sesm.UserId)
		if errors.Is(err, sesm.ErrSessionNotFound) {
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			errInternalSrvResp(err)
			return
		}

		resp, _ := m.users.GetUser(&users.GetUserInput{ID: userId.(int)})
		if resp == nil {
			errInternalSrvResp(err)
			return
		}

		if err := m.checkSessionActivity(r); errors.Is(err, ErrSessionExpired) {
			m.sesManager.DeleteCurrentSession(r)
			http.SetCookie(w, dto.DeleteCookie(sesm.SessionId))
			// flash warning
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			errInternalSrvResp(err)
			return
		}

		if err := m.sessionRateLimiting(r); errors.Is(err, ErrSessionExpired) {
			m.sesManager.DeleteCurrentSession(r)
			http.SetCookie(w, dto.DeleteCookie(sesm.SessionId))
			next.ServeHTTP(w, r)
			return
		} else if errors.Is(err, ErrTooManyRequests) {
			m.exceptions.ErrTooManyRequestsHandler(w, r)
			return
		} else if err != nil {
			errInternalSrvResp(err)
			return
		}

		ctx := context.WithValue(r.Context(), dto.ContextKeyUser, resp.User)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *middleware) checkSessionActivity(r *http.Request) error {
	lastRequest, err := m.sesManager.GetSessionData(r, sesm.LastRequest)
	if err != nil {
		return err
	}

	if time.Now().Sub(lastRequest.(time.Time)) >= sesTimeLimit {
		return ErrSessionExpired
	}

	if err := m.sesManager.UpdateSessionLastRequest(r); err != nil {
		return err
	}

	return nil
}

func (m *middleware) setRateLimiter(sessionId string) {
	m.mutex.RLock()
	limiter, exists := m.limiters[sessionId]
	m.mutex.RUnlock()

	if !exists {
		// Create a new rate limiter for this session
		limiter = make(chan time.Time, maxRequests)
		for i := 0; i < maxRequests; i++ {
			limiter <- time.Now()
		}

		go func() {
			filler := time.NewTicker(200 * time.Millisecond)
			for t := range filler.C {
				select {
				case limiter <- t:
				default: // If the channel is full, discard the tick
				}
			}
		}()

		m.mutex.Lock()
		m.limiters[sessionId] = limiter
		m.mutex.Unlock()
	}
}

func (m *middleware) sessionRateLimiting(r *http.Request) error {
	sessionId, err := m.sesManager.CurrentSessionID(r)
	if err != nil {
		return err
	}

	status, err := m.sesManager.GetSessionData(r, sesm.Status)
	if err != nil {
		return err
	}

	if status.(string) == "blocked" {
		lrVal, err := m.sesManager.GetSessionData(r, sesm.LastRequest)
		if err != nil {
			return err
		}
		lastRequest := lrVal.(time.Time)

		btVal, err := m.sesManager.GetSessionData(r, sesm.BlockTimestamp)
		if err != nil {
			return err
		}
		blockTimestamp := btVal.(time.Time)

		if lastRequest.Sub(blockTimestamp) >= sesBlockTime {
			m.mutex.Lock()
			//delete(m.blockedSessions, sessionId)
			delete(m.limiters, sessionId)
			m.mutex.Unlock()

			return ErrSessionExpired
		} else {
			return ErrTooManyRequests
		}
	}

	m.setRateLimiter(sessionId)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	select {
	case <-m.limiters[sessionId]:
		return nil
	default:
		if err := m.sesManager.BlockSession(r); err != nil {
			return err
		}

		return ErrTooManyRequests
	}

}
