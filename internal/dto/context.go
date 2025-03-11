package dto

import (
	"net/http"
)

type contextKey string

var (
	ContextKeyUser = contextKey("user")
	ContextKeyRole = contextKey("role")
)

func GetAuthUser(r *http.Request) *User {
	val := r.Context().Value(ContextKeyUser)

	user, ok := val.(*User)
	if !ok {
		return nil
	}

	return user
}

func GetUserRole(r *http.Request) string {
	val := r.Context().Value(ContextKeyRole)

	role, ok := val.(string)
	if !ok {
		return ""
	}

	return role
}
