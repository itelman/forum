package users

import (
	"github.com/itelman/forum/internal/service/users/domain"
	"net/mail"
	"regexp"
	"strings"

	"github.com/itelman/forum/pkg/validator"
)

type LoginUserInput struct {
	Username string
	Password string
	Errors   validator.Errors
}

func (i *LoginUserInput) validate() error {
	if !usernameRX.MatchString(i.Username) {
		i.Errors.Add("username", validator.ErrInputRequired("username"))
	}

	if !passwordRX.MatchString(i.Password) {
		i.Errors.Add("password", validator.ErrInputRequired("password"))
	}

	if len(i.Errors) != 0 {
		return domain.ErrUsersBadRequest
	}

	i.Username = strings.ToLower(i.Username)

	return nil
}

type GetUserInput struct {
	ID int
}

var (
	usernameRX = regexp.MustCompile(`^[a-zA-Z]{5,}([._]{0,1}[a-zA-Z0-9]{2,})*$`)
	passwordRX = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
)

const (
	nameMinLen = 5
	nameMaxLen = 30

	pwdMinLen = 6
	pwdMaxLen = 20
)

type SignupUserInput struct {
	Username string
	Email    string
	Password string
	Errors   validator.Errors
}

func (i *SignupUserInput) validate() error {
	i.validateUsername()
	i.validateEmail()
	i.validatePwd()

	if len(i.Errors) != 0 {
		return domain.ErrUsersBadRequest
	}

	i.Username = strings.ToLower(i.Username)
	i.Email = strings.ToLower(i.Email)

	return nil
}

func (i *SignupUserInput) validateUsername() {
	if !usernameRX.MatchString(i.Username) {
		i.Errors.Add("username", validator.ErrInputRequired("username"))
		return
	}

	if !(len(i.Username) >= nameMinLen && len(i.Username) <= nameMaxLen) {
		i.Errors.Add("username", validator.ErrInputLength(nameMinLen, nameMaxLen))
	}
}

func (i *SignupUserInput) validateEmail() {
	email, err := mail.ParseAddress(i.Email)
	if err != nil {
		i.Errors.Add("email", validator.ErrInputRequired("email address"))
		return
	}

	if email.Address != i.Email {
		i.Errors.Add("email", validator.ErrInputRequired("email address"))
	}
}

func (i *SignupUserInput) validatePwd() {
	if !passwordRX.MatchString(i.Password) {
		i.Errors.Add("password", validator.ErrInputRequired("password"))
		return
	}

	if !(len(i.Password) >= pwdMinLen && len(i.Password) <= pwdMaxLen) {
		i.Errors.Add("password", validator.ErrInputLength(pwdMinLen, pwdMaxLen))
	}
}
