package domain

import "errors"

type UsersOAuthRepository interface {
	GetUserID(input GetUserOAuthInput) (int, error)
	Create(input CreateUserOAuthInput) error
}

type GetUserOAuthInput struct {
	AuthTypeID int
	AccountID  string
}

type CreateUserOAuthInput struct {
	UserID     int
	AuthTypeID int
	AccountID  string
}

var (
	ErrOAuthFailed       = errors.New("OAUTH: authentication failed")
	ErrOAuthUserNotFound = errors.New("DATABASE: User (OAUTH) not found")
)
