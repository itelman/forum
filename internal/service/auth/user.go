package auth

import (
	"forum/internal/repository/models"
	"net/http"
)

type contextKey string

var ContextKeyUser = contextKey("user")

func AuthenticatedUser(r *http.Request) *models.User {
	value := r.Context().Value(ContextKeyUser)

	user, ok := value.(*models.User)
	if !ok {
		return nil
	}

	return user
}
