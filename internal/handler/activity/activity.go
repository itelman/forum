package activity

import (
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/handler"
	"github.com/itelman/forum/internal/service/activity"
	"github.com/itelman/forum/pkg/templates"
	"net/http"
)

type handlers struct {
	*handler.Handlers
	activity activity.Service
}

func NewHandlers(handler *handler.Handlers, activity activity.Service) *handlers {
	return &handlers{handler, activity}
}

func (h *handlers) RegisterMux(mux *http.ServeMux) {
	routes := []dto.Route{
		{"/user/activity/created", dto.GetMethod, h.getAllCreatedPosts},
		{"/user/activity/reacted", dto.GetMethod, h.getAllReactedPosts},
		{"/user/activity/commented", dto.GetMethod, h.getAllCommentedPosts},
	}

	for _, route := range routes {
		mux.Handle(route.Path, h.DynMiddleware.Chain(h.DynMiddleware.RequireAuthenticatedUser(http.HandlerFunc(route.Handler)), route.Path, route.Methods))
	}
}

func (h *handlers) getAllCreatedPosts(w http.ResponseWriter, r *http.Request) {
	req := activity.DecodeGetAllCreatedPosts(r)

	resp, err := h.activity.GetAllCreatedPosts(req.(*activity.GetAllCreatedPostsInput))
	if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	if err := h.TmplRender.RenderData(w, r, "activity_page", templates.TemplateData{
		templates.Posts: resp.Posts,
	}); err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}
}

func (h *handlers) getAllReactedPosts(w http.ResponseWriter, r *http.Request) {
	req := activity.DecodeGetAllReactedPosts(r)

	resp, err := h.activity.GetAllReactedPosts(req.(*activity.GetAllReactedPostsInput))
	if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	if err := h.TmplRender.RenderData(w, r, "activity_page", templates.TemplateData{
		templates.Posts: resp.Posts,
	}); err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}
}

func (h *handlers) getAllCommentedPosts(w http.ResponseWriter, r *http.Request) {
	req := activity.DecodeGetAllCommentedPosts(r)

	resp, err := h.activity.GetAllCommentedPosts(req.(*activity.GetAllCommentedPostsInput))
	if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	if err := h.TmplRender.RenderData(w, r, "activity_page", templates.TemplateData{
		templates.Posts: resp.Posts,
	}); err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}
}
