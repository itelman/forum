package users

import (
	"github.com/itelman/forum/internal/service/users/domain"
	"github.com/itelman/forum/pkg/validator"
	"net/http"
)

func DecodeLoginUser(r *http.Request) (interface{}, error) {
	if err := r.ParseForm(); err != nil {
		return nil, domain.ErrUsersBadRequest
	}

	return &LoginUserInput{
		Username: r.PostForm.Get("username"),
		Password: r.PostForm.Get("password"),
		Errors:   make(validator.Errors),
	}, nil
}

func DecodeSignupUser(r *http.Request) (interface{}, error) {
	if err := r.ParseForm(); err != nil {
		return nil, domain.ErrUsersBadRequest
	}

	return &SignupUserInput{
		Username: r.PostForm.Get("username"),
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
		Errors:   make(validator.Errors),
	}, nil
}
