package domain

import (
	"github.com/itelman/forum/internal/dto"
)

type UsersRepository interface {
	Get(input GetUserInput) (*dto.User, error)
}

type GetUserInput struct {
	Key   string
	Value interface{}
}
