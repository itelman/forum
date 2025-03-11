package oauth

type AuthApi interface {
	GetAuthUri() string
	GetUserData(code string) (*AuthData, error)
	GetJSONResponse(code string) ([]byte, error)
}

type AuthData struct {
	AccountID string
	Username  string
	Email     string
	Provider  string
}
