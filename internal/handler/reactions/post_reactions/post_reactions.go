package post_reactions

import (
	"errors"
	"fmt"
	"github.com/itelman/forum/internal/service/post_reactions/domain"
	"net/http"

	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/handler"
	"github.com/itelman/forum/internal/service/post_reactions"
)

type postReactionHandlers struct {
	*handler.Handlers
	postReactions post_reactions.Service
}

func NewHandlers(
	handler *handler.Handlers,
	postReactions post_reactions.Service,
) *postReactionHandlers {
	return &postReactionHandlers{handler, postReactions}
}

func (h *postReactionHandlers) RegisterMux(mux *http.ServeMux) {
	route := dto.Route{Path: "/user/posts/react", Methods: dto.PostMethod, Handler: h.create}
	mux.Handle(route.Path, h.DynMiddleware.Chain(h.DynMiddleware.RequireAuthenticatedUser(http.HandlerFunc(route.Handler)), route.Path, route.Methods))
}

func (h *postReactionHandlers) create(w http.ResponseWriter, r *http.Request) {
	req, err := post_reactions.DecodeCreatePostReaction(r)
	if err != nil {
		h.Exceptions.ErrBadRequestHandler(w, r)
		return
	}

	input := req.(*post_reactions.CreatePostReactionInput)

	if err := h.postReactions.CreatePostReaction(input); errors.Is(err, domain.ErrPostReactionsBadRequest) {
		h.Exceptions.ErrBadRequestHandler(w, r)
		return
	} else if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/posts?id=%d", input.PostID), http.StatusSeeOther)
}
