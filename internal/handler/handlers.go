package handler

import (
	"github.com/itelman/forum/internal/exception"
	"github.com/itelman/forum/internal/middleware/dynamic"
	"github.com/itelman/forum/pkg/sesm"
	"github.com/itelman/forum/pkg/templates"
)

type Handlers struct {
	DynMiddleware dynamic.DynamicMiddleware
	SesManager    sesm.SessionManager
	Exceptions    exception.Exceptions
	TmplRender    templates.TemplateRender
}

func NewHandlers(
	dynMiddleware dynamic.DynamicMiddleware,
	sesManager sesm.SessionManager,
	exceptions exception.Exceptions,
	tmplRender templates.TemplateRender,
) *Handlers {
	return &Handlers{
		DynMiddleware: dynMiddleware,
		SesManager:    sesManager,
		Exceptions:    exceptions,
		TmplRender:    tmplRender,
	}
}
