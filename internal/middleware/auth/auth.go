package auth

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
	sesBlockTime = 3 * time.Hour
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

	limiters     map[int]chan time.Time
	blockedUsers map[int]time.Time
	mutex        sync.RWMutex
}

func NewMiddleware(users users.Service, sesManager sesm.SessionManager, exceptions exception.Exceptions) *middleware {
	return &middleware{
		users:        users,
		sesManager:   sesManager,
		exceptions:   exceptions,
		limiters:     make(map[int]chan time.Time),
		blockedUsers: make(map[int]time.Time),
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

		val, err := m.userRateLimiting(r)
		if errors.Is(err, ErrSessionExpired) {
			m.sesManager.DeleteAllUserSessions(val)
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

func (m *middleware) setRateLimiter(userId int) {
	m.mutex.RLock()
	limiter, exists := m.limiters[userId]
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
		m.limiters[userId] = limiter
		m.mutex.Unlock()
	}
}

func (m *middleware) userRateLimiting(r *http.Request) (int, error) {
	uidVal, err := m.sesManager.GetSessionData(r, sesm.UserId)
	if err != nil {
		return -1, err
	}
	userId := uidVal.(int)

	lrVal, err := m.sesManager.GetSessionData(r, sesm.LastRequest)
	if err != nil {
		return -1, err
	}
	lastRequest := lrVal.(time.Time)

	m.mutex.RLock()
	blockTimestamp, exists := m.blockedUsers[userId]
	m.mutex.RUnlock()

	if exists {
		if lastRequest.Sub(blockTimestamp) >= sesBlockTime {
			m.mutex.Lock()
			delete(m.blockedUsers, userId)
			delete(m.limiters, userId)
			m.mutex.Unlock()

			return userId, ErrSessionExpired
		} else {
			return -1, ErrTooManyRequests
		}
	}

	m.setRateLimiter(userId)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	select {
	case <-m.limiters[userId]:
		return -1, nil
	default:
		m.blockedUsers[userId] = time.Now()
		return -1, ErrTooManyRequests
	}

}
