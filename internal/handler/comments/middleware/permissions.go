package middleware

import (
	"context"
	"errors"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/exception"
	"github.com/itelman/forum/internal/service/comments"
	"github.com/itelman/forum/internal/service/comments/domain"
	"net/http"
)

type contextKey string

var (
	ContextKeyComment = contextKey("comment")
)

type CommentsCheckPermissionMiddleware interface {
	CheckUserPermissions(next http.Handler) http.Handler
}

type middleware struct {
	comments   comments.Service
	exceptions exception.Exceptions
}

func NewMiddleware(comments comments.Service, exceptions exception.Exceptions) *middleware {
	return &middleware{
		comments:   comments,
		exceptions: exceptions,
	}
}

func (m *middleware) CheckUserPermissions(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		commReq, err := comments.DecodeGetComment(r)
		if err != nil {
			m.exceptions.ErrBadRequestHandler(w, r)
			return
		}

		commResp, err := m.comments.GetComment(commReq.(*comments.GetCommentInput))
		if errors.Is(err, domain.ErrCommentNotFound) {
			m.exceptions.ErrNotFoundHandler(w, r)
			return
		} else if err != nil {
			m.exceptions.ErrInternalServerHandler(w, r, err)
			return
		}

		if commResp.Comment.User.ID != dto.GetAuthUser(r).ID {
			m.exceptions.ErrForbiddenHandler(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyComment, commResp.Comment)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetCommentFromContext(r *http.Request) *dto.Comment {
	val := r.Context().Value(ContextKeyComment)

	comment, ok := val.(*dto.Comment)
	if !ok {
		return nil
	}

	return comment
}
