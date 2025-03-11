package post_reactions

import (
	"errors"
	"fmt"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/handler"
	"github.com/itelman/forum/internal/service/comments/domain"
	"github.com/itelman/forum/internal/service/post_reactions"
	"github.com/itelman/forum/internal/service/posts"
	"net/http"
)

type postReactionHandlers struct {
	*handler.Handlers
	postReactions post_reactions.Service
	posts         posts.Service
}

func NewHandlers(
	handler *handler.Handlers,
	postReactions post_reactions.Service,
	posts posts.Service,
) *postReactionHandlers {
	return &postReactionHandlers{handler, postReactions, posts}
}

func (h *postReactionHandlers) RegisterMux(mux *http.ServeMux) {
	route := dto.Route{"/user/posts/react", dto.PostMethod, h.create}
	mux.Handle(route.Path, h.DynMiddleware.Chain(h.DynMiddleware.RequireAuthenticatedUser(http.HandlerFunc(route.Handler)), route.Path, route.Methods))
}

func (h *postReactionHandlers) create(w http.ResponseWriter, r *http.Request) {
	req, err := post_reactions.DecodeCreatePostReaction(r)

	input := req.(*post_reactions.CreatePostReactionInput)

	_, err = h.posts.GetPost(&posts.GetPostInput{ID: input.PostID})
	if errors.Is(err, domain.ErrCommentNotFound) {
		h.Exceptions.ErrNotFoundHandler(w, r)
		return
	} else if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	if err := h.postReactions.CreatePostReaction(input); err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/posts?id=%d", input.PostID), http.StatusSeeOther)
}
