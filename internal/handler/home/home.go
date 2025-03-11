package home

import (
	"encoding/json"
	"errors"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/handler"
	"github.com/itelman/forum/internal/service/categories"
	"github.com/itelman/forum/internal/service/filters"
	"github.com/itelman/forum/internal/service/filters/domain"
	"github.com/itelman/forum/internal/service/posts"
	"github.com/itelman/forum/pkg/templates"
	"github.com/itelman/forum/pkg/validator"
	"net/http"
)

type handlers struct {
	*handler.Handlers
	posts      posts.Service
	categories categories.Service
	filters    filters.Service
}

func NewHandlers(handler *handler.Handlers, posts posts.Service, categories categories.Service, filters filters.Service) *handlers {
	return &handlers{handler, posts, categories, filters}
}

func (h *handlers) RegisterMux(mux *http.ServeMux) {
	routes := []dto.Route{
		{"/", dto.GetMethod, h.home},
		{"/results", dto.PostMethod, h.results},
		{"/health", dto.GetMethod, h.healthCheck},
	}

	for _, route := range routes {
		mux.Handle(route.Path, h.DynMiddleware.Chain(http.HandlerFunc(route.Handler), route.Path, route.Methods))
	}
}

func (h *handlers) healthCheck(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"status": "available",
		"system_info": map[string]string{
			"environment": "development",
			"version":     "1.0",
		},
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		h.Exceptions.ErrBadRequestHandler(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}

func (h *handlers) home(w http.ResponseWriter, r *http.Request) {
	postsResp, err := h.posts.GetAllLatestPosts()
	if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	catgRsp, err := h.categories.GetAllCategories()
	if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	if err := h.TmplRender.RenderData(w, r, "home_page", templates.TemplateData{
		templates.Posts:      postsResp.Posts,
		templates.Categories: catgRsp.Categories,
		templates.Form:       validator.NewForm(nil, nil),
	}); err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}
}

func (h *handlers) results(w http.ResponseWriter, r *http.Request) {
	req, err := filters.DecodeGetPostsByFilters(r)
	if err != nil {
		h.Exceptions.ErrBadRequestHandler(w, r)
		return
	}

	input := req.(*filters.GetPostsByFiltersInput)

	filtersResp, err := h.filters.GetPostsByFilters(input)
	if errors.Is(err, domain.ErrFiltersNoneSelected) {
		if dto.GetAuthUser(r) != nil {
			if err := h.SesManager.UpdateSessionFlash(r, dto.FlashFilterSelect); err != nil {
				h.Exceptions.ErrInternalServerHandler(w, r, err)
				return
			}
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else if errors.Is(err, domain.ErrUserUnauthorized) {
		http.Redirect(w, r, "/user/login", http.StatusFound)
		return
	} else if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	catgRsp, err := h.categories.GetAllCategories()
	if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	if err := h.TmplRender.RenderData(w, r, "home_page", templates.TemplateData{
		templates.Posts:      filtersResp.Posts,
		templates.Categories: catgRsp.Categories,
		templates.Form:       validator.NewForm(nil, input.Errors),
	}); err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}
}
