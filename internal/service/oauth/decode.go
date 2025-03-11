package oauth

import (
	"github.com/itelman/forum/internal/service/oauth/domain"
	"github.com/itelman/forum/pkg/oauth"
	"net/http"
)

func DecodeLoginUserInput(r *http.Request, api oauth.AuthApi) (interface{}, error) {
	code, exists := r.URL.Query()["code"]
	if !exists {
		return nil, domain.ErrOAuthFailed
	}

	authData, err := api.GetUserData(code[0])
	if err != nil {
		return nil, err
	}

	return &LoginUserInput{
		AccountID: authData.AccountID,
		Username:  authData.Username,
		Email:     authData.Email,
		Provider:  authData.Provider,
	}, nil
}
