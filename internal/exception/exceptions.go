package exception

import (
	"fmt"
	"github.com/itelman/forum/pkg/templates"
	"log"
	"net/http"
	"runtime"
	"runtime/debug"
)

type Exceptions interface {
	ErrBadRequestHandler(w http.ResponseWriter, r *http.Request)
	ErrUnauthorizedHandler(w http.ResponseWriter, r *http.Request)
	ErrForbiddenHandler(w http.ResponseWriter, r *http.Request)
	ErrNotFoundHandler(w http.ResponseWriter, r *http.Request)
	ErrNotAllowedHandler(w http.ResponseWriter, r *http.Request)
	ErrTooManyRequestsHandler(w http.ResponseWriter, r *http.Request)
	ErrInternalServerHandler(w http.ResponseWriter, r *http.Request, err error)
}

type exceptions struct {
	errorLog   *log.Logger
	tmplRender templates.TemplateRender
}

func NewExceptions(errorLog *log.Logger, tmplRender templates.TemplateRender) *exceptions {
	return &exceptions{errorLog: errorLog, tmplRender: tmplRender}
}

func (e *exceptions) defaultErrorHandler(w http.ResponseWriter, r *http.Request, data errorData) {
	w.WriteHeader(data.Code)
	if err := e.tmplRender.RenderData(w, r, "error_page", templates.TemplateData{templates.Error: data}); err != nil {
		e.errInternalServerText(w)
		return
	}
}

func (e *exceptions) ErrBadRequestHandler(w http.ResponseWriter, r *http.Request) {
	e.defaultErrorHandler(w, r, errBadRequestData)
	return
}

func (e *exceptions) ErrUnauthorizedHandler(w http.ResponseWriter, r *http.Request) {
	e.defaultErrorHandler(w, r, errUnauthorizedData)
	return
}

func (e *exceptions) ErrForbiddenHandler(w http.ResponseWriter, r *http.Request) {
	e.defaultErrorHandler(w, r, errForbiddenData)
	return
}

func (e *exceptions) ErrNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	e.defaultErrorHandler(w, r, errNotFoundData)
	return
}

func (e *exceptions) ErrNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	e.defaultErrorHandler(w, r, errNotAllowedData)
	return
}

func (e *exceptions) ErrTooManyRequestsHandler(w http.ResponseWriter, r *http.Request) {
	e.defaultErrorHandler(w, r, errTooManyRequestsData)
	return
}

func (e *exceptions) ErrInternalServerHandler(w http.ResponseWriter, r *http.Request, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	e.errorLog.Output(2, trace)
	//e.errorLog.Printf("ERROR: %v (caller: %s)", err, getCaller())

	e.defaultErrorHandler(w, r, errInternalServerData)
	return
}

func (e *exceptions) errInternalServerText(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	return
}

func getCaller() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown"
	}
	return fmt.Sprintf("%s:%d", file, line)
}
