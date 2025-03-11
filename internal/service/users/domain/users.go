package domain

import (
	"errors"
	"github.com/itelman/forum/internal/dto"
)

type UsersRepository interface {
	Create(input RegisterUserInput) error
	Get(input GetUserInput) (*dto.User, error)
	Authenticate(input AuthUserInput) (int, error)
}

type GetUserInput struct {
	Key   string
	Value interface{}
}

type AuthUserInput struct {
	Username string
	Password string
}

type RegisterUserInput struct {
	Username string
	Email    string
	Password string
}

var (
	ErrUsersBadRequest    = errors.New("USERS: bad request")
	ErrUserNotFound       = errors.New("DATABASE: User not found")
	ErrUserExists         = errors.New("DATABASE: User exists")
	ErrInvalidCredentials = errors.New("DATABASE: Invalid credentials")
)
