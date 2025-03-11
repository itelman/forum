package oauth

import (
	"database/sql"
	"errors"
	"github.com/itelman/forum/internal/service/oauth/adapters"
	"github.com/itelman/forum/internal/service/oauth/domain"
	"strings"
)

type Service interface {
	LoginUser(input *LoginUserInput) (*LoginUserResponse, error)
}

type service struct {
	users      domain.UsersRepository
	usersOAuth domain.UsersOAuthRepository
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
		s.usersOAuth = adapters.NewUsersOAuthRepositorySqlite(db)
	}
}

type GetAuthCodeInput struct {
	Code string
}

type LoginUserInput struct {
	AccountID string
	Username  string
	Email     string
	Provider  string
}

func (i *LoginUserInput) getProviderID() int {
	switch strings.ToLower(i.Provider) {
	case "github":
		return 1
	case "google":
		return 2
	default:
		return -1
	}
}

type LoginUserResponse struct {
	UserID int
}

func (s *service) LoginUser(input *LoginUserInput) (*LoginUserResponse, error) {
	if input.getProviderID() == -1 {
		return nil, domain.ErrOAuthFailed
	}

	userId, err := s.usersOAuth.GetUserID(domain.GetUserOAuthInput{
		AuthTypeID: input.getProviderID(),
		AccountID:  input.AccountID,
	})
	if err != nil && !errors.Is(err, domain.ErrOAuthUserNotFound) {
		return nil, err
	}

	if userId != -1 {
		return &LoginUserResponse{UserID: userId}, nil
	}

	user, err := s.users.Get(domain.GetUserInput{
		Key:   "email",
		Value: input.Email,
	})
	if err != nil {
		return nil, err
	}

	if err := s.usersOAuth.Create(domain.CreateUserOAuthInput{
		UserID:     user.ID,
		AuthTypeID: input.getProviderID(),
		AccountID:  input.AccountID,
	}); err != nil {
		return nil, err
	}

	return &LoginUserResponse{UserID: user.ID}, nil
}
