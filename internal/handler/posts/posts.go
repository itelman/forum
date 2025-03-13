package posts

import (
	"errors"
	"fmt"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/handler"
	"github.com/itelman/forum/internal/handler/posts/middleware"
	"github.com/itelman/forum/internal/service/categories"
	"github.com/itelman/forum/internal/service/posts"
	"github.com/itelman/forum/internal/service/posts/domain"
	"github.com/itelman/forum/pkg/templates"
	"github.com/itelman/forum/pkg/validator"
	"net/http"
	"net/url"
)

type handlers struct {
	*handler.Handlers
	checkPerm     middleware.PostsCheckPermissionMiddleware
	posts         posts.Service
	categories    categories.Service
	postImagesDir string
}

func NewHandlers(handler *handler.Handlers, posts posts.Service, categories categories.Service, dir string) *handlers {
	checkPermMid := middleware.NewMiddleware(posts, handler.Exceptions)
	return &handlers{handler, checkPermMid, posts, categories, dir}
}

func (h *handlers) RegisterMux(mux *http.ServeMux) {
	showRoute := dto.Route{Path: "/posts", Methods: dto.GetMethod, Handler: h.get}
	mux.Handle(showRoute.Path, h.DynMiddleware.Chain(http.HandlerFunc(showRoute.Handler), showRoute.Path, showRoute.Methods))

	createRoute := dto.Route{Path: "/user/posts/create", Methods: dto.GetPostMethods, Handler: h.createForm}
	mux.Handle(createRoute.Path, h.DynMiddleware.Chain(h.DynMiddleware.RequireAuthenticatedUser(http.HandlerFunc(createRoute.Handler)), createRoute.Path, createRoute.Methods))

	editDeleteRoutes := []dto.Route{
		{"/user/posts/edit", dto.GetPostMethods, h.editForm},
		{"/user/posts/delete", dto.GetMethod, h.delete},
	}

	for _, route := range editDeleteRoutes {
		mux.Handle(route.Path, h.DynMiddleware.Chain(h.DynMiddleware.RequireAuthenticatedUser(h.checkPerm.CheckUserPermissions(http.HandlerFunc(route.Handler))), route.Path, route.Methods))
	}
}

func (h *handlers) createForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.create(w, r)
		return
	}

	catgRsp, err := h.categories.GetAllCategories()
	if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	if err := h.TmplRender.RenderData(w, r, "create_page", templates.TemplateData{
		templates.Form:       validator.NewForm(nil, nil),
		templates.Categories: catgRsp.Categories,
	}); err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}
}

func (h *handlers) create(w http.ResponseWriter, r *http.Request) {
	req, err := posts.DecodeCreatePost(r)
	if errors.Is(err, domain.ErrPostsBadRequest) {
		h.Exceptions.ErrBadRequestHandler(w, r)
		return
	} else if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	input := req.(*posts.CreatePostInput)

	resp, err := h.posts.CreatePost(input, h.postImagesDir)
	if errors.Is(err, domain.ErrPostsBadRequest) {
		catgRsp, err := h.categories.GetAllCategories()
		if err != nil {
			h.Exceptions.ErrInternalServerHandler(w, r, err)
			return
		}

		if err := h.TmplRender.RenderData(w, r, "create_page", templates.TemplateData{
			templates.Form:       validator.NewForm(r.PostForm, input.Errors),
			templates.Categories: catgRsp.Categories,
		}); err != nil {
			h.Exceptions.ErrInternalServerHandler(w, r, err)
			return
		}

		return
	} else if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/posts?id=%d", resp.PostID), http.StatusSeeOther)
}

func (h *handlers) get(w http.ResponseWriter, r *http.Request) {
	postReq, err := posts.DecodeGetPost(r)
	if err != nil {
		h.Exceptions.ErrBadRequestHandler(w, r)
		return
	}

	postResp, err := h.posts.GetPost(postReq.(*posts.GetPostInput))
	if errors.Is(err, domain.ErrPostNotFound) {
		h.Exceptions.ErrNotFoundHandler(w, r)
		return
	} else if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	if err := h.TmplRender.RenderData(w, r, "show_page", templates.TemplateData{
		templates.Post: postResp.Post,
	}); err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}
}

func (h *handlers) delete(w http.ResponseWriter, r *http.Request) {
	req, err := posts.DecodeDeletePost(r)
	if err != nil {
		h.Exceptions.ErrBadRequestHandler(w, r)
		return
	}

	input := req.(*posts.DeletePostInput)

	if err := h.posts.DeletePost(input, h.postImagesDir); errors.Is(err, domain.ErrPostNotFound) {
		h.Exceptions.ErrNotFoundHandler(w, r)
		return
	} else if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	if err := h.SesManager.UpdateSessionFlash(r, dto.FlashPostRemoved); err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *handlers) editForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.edit(w, r)
		return
	}

	post := middleware.GetPostFromContext(r)

	autoForm := make(url.Values)
	autoForm.Set("title", post.Title)
	autoForm.Set("content", post.Content)

	if err := h.TmplRender.RenderData(w, r, "edit_post_page", templates.TemplateData{
		templates.Post: post,
		templates.Form: validator.NewForm(autoForm, nil),
	}); err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}
}

func (h *handlers) edit(w http.ResponseWriter, r *http.Request) {
	req, err := posts.DecodeUpdatePost(r)
	if err != nil {
		h.Exceptions.ErrBadRequestHandler(w, r)
		return
	}

	input := req.(*posts.UpdatePostInput)
	post := middleware.GetPostFromContext(r)

	if err := h.posts.UpdatePost(input, post); errors.Is(err, domain.ErrPostsBadRequest) {
		if err := h.TmplRender.RenderData(w, r, "edit_post_page", templates.TemplateData{
			templates.Post: post,
			templates.Form: validator.NewForm(r.PostForm, input.Errors),
		}); err != nil {
			h.Exceptions.ErrInternalServerHandler(w, r, err)
			return
		}

		return
	} else if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/posts?id=%d", post.ID), http.StatusSeeOther)
}
