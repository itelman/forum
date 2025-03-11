package notifications

import (
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/handler"
	"github.com/itelman/forum/internal/service/notifications"
	"github.com/itelman/forum/pkg/templates"
	"github.com/itelman/forum/pkg/validator"
	"net/http"
)

type handlers struct {
	*handler.Handlers
	notifications notifications.Service
}

func NewHandlers(handler *handler.Handlers, notifications notifications.Service) *handlers {
	return &handlers{handler, notifications}
}

func (h *handlers) RegisterMux(mux *http.ServeMux) {
	routes := []dto.Route{
		{Path: "/user/notifications", Methods: dto.GetPostMethods, Handler: h.redirect},
		{Path: "/user/notifications/comments", Methods: dto.GetMethod, Handler: h.getCommentNotifications},
		{Path: "/user/notifications/reactions", Methods: dto.GetMethod, Handler: h.getPostReactionNotifications},
	}

	for _, route := range routes {
		mux.Handle(route.Path, h.DynMiddleware.Chain(h.DynMiddleware.RequireAuthenticatedUser(http.HandlerFunc(route.Handler)), route.Path, route.Methods))
	}
}

func (h *handlers) redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/user/notifications/comments", http.StatusMovedPermanently)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.Exceptions.ErrBadRequestHandler(w, r)
		return
	}

	if r.PostForm.Get("filter") == "1" {
		http.Redirect(w, r, "/user/notifications/comments", http.StatusMovedPermanently)
		return
	} else if r.PostForm.Get("filter") == "2" {
		http.Redirect(w, r, "/user/notifications/reactions", http.StatusMovedPermanently)
		return
	}

	if err := h.SesManager.UpdateSessionFlash(r, dto.FlashFilterSelect); err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}
	http.Redirect(w, r, "/user/notifications/comments", http.StatusMovedPermanently)
}

func (h *handlers) getCommentNotifications(w http.ResponseWriter, r *http.Request) {
	resp, err := h.notifications.GetAllCommentNotifications(notifications.DecodeGetAllCommentNotifications(r).(*notifications.GetAllCommentNotificationsInput))
	if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	if err := h.TmplRender.RenderData(w, r, "notifications_page", templates.TemplateData{
		templates.Comments: resp.Comments,
		templates.Form:     validator.NewForm(nil, nil),
	}); err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}
}

func (h *handlers) getPostReactionNotifications(w http.ResponseWriter, r *http.Request) {
	resp, err := h.notifications.GetAllPostReactionNotifications(notifications.DecodeGetAllPostReactionNotifications(r).(*notifications.GetAllPostReactionNotificationsInput))
	if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	if err := h.TmplRender.RenderData(w, r, "notifications_page", templates.TemplateData{
		templates.PostReactions: resp.PostReactions,
		templates.Form:          validator.NewForm(nil, nil),
	}); err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}
}
