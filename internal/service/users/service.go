package users

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/itelman/forum/internal/service/users/domain"

	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/users/adapters"
)

type Service interface {
	SignupUser(input *SignupUserInput) error
	LoginUser(input *LoginUserInput) (*LoginUserResponse, error)
	GetUser(input *GetUserInput) (*GetUserResponse, error)
}

type service struct {
	users domain.UsersRepository
}

func NewService(opts ...Option) *service {
	svc := &service{}
	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

type Option func(*service)

func WithSqlite(db *sql.DB) Option {
	return func(s *service) {
		s.users = adapters.NewUsersRepositorySqlite(db)
	}
}

func (s *service) SignupUser(input *SignupUserInput) error {
	if err := input.validate(); err != nil {
		return err
	}

	inputSet := []domain.GetUserInput{
		{Key: "username", Value: input.Username},
		{Key: "email", Value: input.Email},
	}

	for _, getUserInput := range inputSet {
		user, err := s.users.Get(getUserInput)
		if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
			return err
		}

		if user != nil {
			input.Errors.Add(getUserInput.Key, fmt.Sprintf("An account with such %s already exists", getUserInput.Key))
		}
	}

	if len(input.Errors) != 0 {
		return domain.ErrUserExists
	}

	if err := s.users.Create(domain.RegisterUserInput{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}); err != nil {
		return err
	}

	return nil
}

type LoginUserResponse struct {
	UserID int
}

func (s *service) LoginUser(input *LoginUserInput) (*LoginUserResponse, error) {
	if err := input.validate(); err != nil {
		return nil, err
	}

	userId, err := s.users.Authenticate(domain.AuthUserInput{
		Username: input.Username,
		Password: input.Password,
	})
	if errors.Is(err, domain.ErrUserNotFound) {
		input.Errors.Add("username", "No account found with such username")
	} else if errors.Is(err, domain.ErrInvalidCredentials) {
		input.Errors.Add("generic", "Authentication failed. Please check your credentials and try again")
	} else if err != nil {
		return nil, err
	}

	return &LoginUserResponse{userId}, nil
}

type GetUserResponse struct {
	User *dto.User
}

func (s *service) GetUser(input *GetUserInput) (*GetUserResponse, error) {
	user, err := s.users.Get(domain.GetUserInput{
		Key:   "id",
		Value: input.ID,
	})
	if err != nil {
		return nil, err
	}

	return &GetUserResponse{user}, nil
}
