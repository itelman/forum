package comments

import (
	"errors"
	"fmt"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/handler"
	"github.com/itelman/forum/internal/handler/comments/middleware"
	"github.com/itelman/forum/internal/service/comments"
	"github.com/itelman/forum/internal/service/comments/domain"
	postDomain "github.com/itelman/forum/internal/service/posts/domain"
	"github.com/itelman/forum/pkg/templates"
	"github.com/itelman/forum/pkg/validator"
	"net/http"
	"net/url"
)

type handlers struct {
	*handler.Handlers
	checkPerm middleware.CommentsCheckPermissionMiddleware
	comments  comments.Service
}

func NewHandlers(handler *handler.Handlers, comments comments.Service) *handlers {
	checkPermMid := middleware.NewMiddleware(comments, handler.Exceptions)
	return &handlers{handler, checkPermMid, comments}
}

func (h *handlers) RegisterMux(mux *http.ServeMux) {
	createRoute := dto.Route{Path: "/user/posts/comments/create", Methods: dto.PostMethod, Handler: h.create}
	mux.Handle(createRoute.Path, h.DynMiddleware.Chain(h.DynMiddleware.RequireAuthenticatedUser(http.HandlerFunc(createRoute.Handler)), createRoute.Path, createRoute.Methods))

	editDeleteRoutes := []dto.Route{
		{Path: "/user/posts/comments/edit", Methods: dto.GetPostMethods, Handler: h.editForm},
		{Path: "/user/posts/comments/delete", Methods: dto.GetMethod, Handler: h.delete},
	}

	for _, route := range editDeleteRoutes {
		mux.Handle(route.Path, h.DynMiddleware.Chain(h.DynMiddleware.RequireAuthenticatedUser(h.checkPerm.CheckUserPermissions(http.HandlerFunc(route.Handler))), route.Path, route.Methods))
	}
}

func (h *handlers) create(w http.ResponseWriter, r *http.Request) {
	req, err := comments.DecodeCreateComment(r)
	if err != nil {
		h.Exceptions.ErrBadRequestHandler(w, r)
		return
	}

	input := req.(*comments.CreateCommentInput)

	if err = h.comments.CreateComment(input); errors.Is(err, domain.ErrCommentsBadRequest) {
		if err := h.SesManager.UpdateSessionFlash(r, dto.FlashCommentEnter); err != nil {
			h.Exceptions.ErrInternalServerHandler(w, r, err)
			return
		}
	} else if errors.Is(err, postDomain.ErrPostNotFound) {
		h.Exceptions.ErrNotFoundHandler(w, r)
		return
	} else if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/posts?id=%d", input.PostID), http.StatusSeeOther)
}

func (h *handlers) delete(w http.ResponseWriter, r *http.Request) {
	req, err := comments.DecodeDeleteComment(r)
	if err != nil {
		h.Exceptions.ErrBadRequestHandler(w, r)
		return
	}

	comment := middleware.GetCommentFromContext(r)

	if err = h.comments.DeleteComment(req.(*comments.DeleteCommentInput)); err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/posts?id=%d", comment.PostID), http.StatusSeeOther)
}

func (h *handlers) editForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.edit(w, r)
		return
	}

	comment := middleware.GetCommentFromContext(r)

	autoForm := make(url.Values)
	autoForm.Set("content", comment.Content)

	if err := h.TmplRender.RenderData(w, r, "edit_comment_page", templates.TemplateData{
		templates.Comment: comment,
		templates.Form:    validator.NewForm(autoForm, nil),
	}); err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}
}

func (h *handlers) edit(w http.ResponseWriter, r *http.Request) {
	req, err := comments.DecodeUpdateComment(r)
	if err != nil {
		h.Exceptions.ErrBadRequestHandler(w, r)
		return
	}

	input := req.(*comments.UpdateCommentInput)
	comment := middleware.GetCommentFromContext(r)

	if err := h.comments.UpdateComment(input, comment); errors.Is(err, domain.ErrCommentsBadRequest) {
		if err := h.TmplRender.RenderData(w, r, "edit_comment_page", templates.TemplateData{
			templates.Comment: comment,
			templates.Form:    validator.NewForm(r.PostForm, input.Errors),
		}); err != nil {
			h.Exceptions.ErrInternalServerHandler(w, r, err)
			return
		}

		return
	} else if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/posts?id=%d", comment.PostID), http.StatusSeeOther)
}
