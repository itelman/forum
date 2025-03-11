package dynamic

import (
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/exception"
	"github.com/itelman/forum/internal/middleware/auth"
	"github.com/itelman/forum/pkg/sesm"
	"net/http"
)

type DynamicMiddleware interface {
	Chain(next http.Handler, path string, methods []string) http.Handler
	RequireAuthenticatedUser(next http.Handler) http.Handler
	ForbidAuthenticatedUser(next http.Handler) http.Handler
	RoleAccessControl(next http.Handler, role string) http.Handler
}

type middleware struct {
	sesManager sesm.SessionManager
	exceptions exception.Exceptions
	authMid    auth.AuthMiddleware
}

func NewMiddleware(authMid auth.AuthMiddleware, sesManager sesm.SessionManager, exceptions exception.Exceptions) *middleware {
	return &middleware{
		authMid:    authMid,
		sesManager: sesManager,
		exceptions: exceptions,
	}
}

func (m *middleware) Chain(next http.Handler, path string, methods []string) http.Handler {
	return m.requestValidation(m.authMid.Authenticate(next), path, methods)
}

func (m *middleware) requestValidation(next http.Handler, path string, methods []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			m.exceptions.ErrNotFoundHandler(w, r)
			return
		}

		for _, method := range methods {
			if method == r.Method {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("Allow", methods[0])
		m.exceptions.ErrNotAllowedHandler(w, r)
	})
}

func (m *middleware) RequireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if dto.GetAuthUser(r) == nil {
			http.Redirect(w, r, "/user/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *middleware) ForbidAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if dto.GetAuthUser(r) != nil {
			m.exceptions.ErrForbiddenHandler(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *middleware) RoleAccessControl(next http.Handler, role string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if dto.GetUserRole(r) != role {
			m.exceptions.ErrForbiddenHandler(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
