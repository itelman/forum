package comment_reactions

import (
	"errors"
	"fmt"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/handler"
	"github.com/itelman/forum/internal/service/comment_reactions"
	"github.com/itelman/forum/internal/service/comments"
	"github.com/itelman/forum/internal/service/comments/domain"
	"net/http"
)

type commentReactionHandlers struct {
	*handler.Handlers
	commentReactions comment_reactions.Service
	comments         comments.Service
}

func NewHandlers(
	handler *handler.Handlers,
	commentReactions comment_reactions.Service,
	comments comments.Service,
) *commentReactionHandlers {
	return &commentReactionHandlers{handler, commentReactions, comments}
}

func (h *commentReactionHandlers) RegisterMux(mux *http.ServeMux) {
	route := dto.Route{"/user/posts/comments/react", dto.PostMethod, h.create}
	mux.Handle(route.Path, h.DynMiddleware.Chain(h.DynMiddleware.RequireAuthenticatedUser(http.HandlerFunc(route.Handler)), route.Path, route.Methods))
}

func (h *commentReactionHandlers) create(w http.ResponseWriter, r *http.Request) {
	req, err := comment_reactions.DecodeCreateCommentReaction(r)

	input := req.(*comment_reactions.CreateCommentReactionInput)

	commResp, err := h.comments.GetComment(&comments.GetCommentInput{ID: input.CommentID})
	if errors.Is(err, domain.ErrCommentNotFound) {
		h.Exceptions.ErrNotFoundHandler(w, r)
		return
	} else if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	if err := h.commentReactions.CreateCommentReaction(input); err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/posts?id=%d", commResp.Comment.PostID), http.StatusSeeOther)
}
