package standard

import (
	"fmt"
	"github.com/itelman/forum/internal/exception"
	"log"
	"net/http"
)

type StandardMiddleware interface {
	Chain(next http.Handler) http.Handler
}

type middleware struct {
	exceptions exception.Exceptions
	infoLog    *log.Logger
}

func NewMiddleware(exceptions exception.Exceptions, infoLog *log.Logger) *middleware {
	return &middleware{exceptions: exceptions, infoLog: infoLog}
}

func (m *middleware) Chain(next http.Handler) http.Handler {
	return m.recoverPanic(m.requestLogging(m.secureHeaders(next)))
}

func (m *middleware) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				m.exceptions.ErrInternalServerHandler(w, r, fmt.Errorf("%s", err))
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (m *middleware) requestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (m *middleware) secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}
