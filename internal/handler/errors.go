package handler

import (
	"fmt"
	"forum/internal/repository/models"
	"forum/internal/service/tmpldata"
	"net/http"
	"runtime/debug"
)

func (h *Handlers) ServerErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	h.App.ErrorLog.Output(2, trace)

	errorModel := &models.Error{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
	}

	w.WriteHeader(http.StatusInternalServerError)
	err = h.App.Render(w, r, &tmpldata.TemplateData{
		TemplateName: "error_page.html",
		Error:        errorModel,
	})
	if err != nil {
		h.ServerErrorTxt(w)
		return
	}
}

func (h *Handlers) ClientErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	errorModel := &models.Error{
		Code:    status,
		Message: http.StatusText(status),
	}

	w.WriteHeader(status)
	err := h.App.Render(w, r, &tmpldata.TemplateData{
		TemplateName: "error_page.html",
		Error:        errorModel,
	})
	if err != nil {
		h.ServerErrorTxt(w)
		return
	}
}

func (h *Handlers) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	h.ClientErrorHandler(w, r, http.StatusNotFound)
	return
}

func (h *Handlers) ServerErrorTxt(w http.ResponseWriter) {
	w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
	return
}
