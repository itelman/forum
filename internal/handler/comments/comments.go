package comments

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/handler"
	"github.com/itelman/forum/internal/service/comments"
	"github.com/itelman/forum/internal/service/comments/domain"
	postDomain "github.com/itelman/forum/internal/service/posts/domain"
	"github.com/itelman/forum/pkg/templates"
	"github.com/itelman/forum/pkg/validator"
)

type handlers struct {
	*handler.Handlers
	comments comments.Service
}

func NewHandlers(handler *handler.Handlers, comments comments.Service) *handlers {
	return &handlers{handler, comments}
}

func (h *handlers) RegisterMux(mux *http.ServeMux) {
	routes := []dto.Route{
		{Path: "/user/posts/comments/create", Methods: dto.PostMethod, Handler: h.create},
		{Path: "/user/posts/comments/edit", Methods: dto.GetPostMethods, Handler: h.editForm},
		{Path: "/user/posts/comments/delete", Methods: dto.GetMethod, Handler: h.delete},
	}

	for _, route := range routes {
		mux.Handle(route.Path, h.DynMiddleware.Chain(h.DynMiddleware.RequireAuthenticatedUser(http.HandlerFunc(route.Handler)), route.Path, route.Methods))
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

	input := req.(*comments.DeleteCommentInput)

	resp, err := h.comments.GetComment(&comments.GetCommentInput{ID: input.ID})
	if errors.Is(err, domain.ErrCommentNotFound) {
		h.Exceptions.ErrNotFoundHandler(w, r)
		return
	}

	if err = h.comments.DeleteComment(input); err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/posts?id=%d", resp.Comment.PostID), http.StatusSeeOther)
}

func (h *handlers) editForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.edit(w, r)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		h.Exceptions.ErrBadRequestHandler(w, r)
		return
	}

	resp, err := h.comments.GetComment(&comments.GetCommentInput{ID: id})
	if errors.Is(err, domain.ErrCommentNotFound) {
		h.Exceptions.ErrNotFoundHandler(w, r)
		return
	} else if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	autoForm := make(url.Values)
	autoForm.Set("content", resp.Comment.Content)

	if err := h.TmplRender.RenderData(w, r, "edit_comment_page", templates.TemplateData{
		templates.Comment: resp.Comment,
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

	commResp, err := h.comments.GetComment(&comments.GetCommentInput{ID: input.ID})
	if errors.Is(err, domain.ErrCommentNotFound) {
		h.Exceptions.ErrNotFoundHandler(w, r)
		return
	} else if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	input.Comment = commResp.Comment

	if err := h.comments.UpdateComment(input); errors.Is(err, domain.ErrCommentsBadRequest) {
		if err := h.TmplRender.RenderData(w, r, "edit_comment_page", templates.TemplateData{
			templates.Comment: commResp.Comment,
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

	http.Redirect(w, r, fmt.Sprintf("/posts?id=%d", commResp.Comment.PostID), http.StatusSeeOther)
}
