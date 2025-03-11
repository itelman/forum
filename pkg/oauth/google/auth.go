package google

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/itelman/forum/pkg/oauth"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type refreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	Expires      int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

type idTokenResponse struct {
	AccessToken string `json:"access_token"`
	Expires     int    `json:"expires_in"`
	IdToken     string `json:"id_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type oAuthApi struct {
	clientSecret string
	clientID     string
	authUri      string
	callbackUri  string
}

func NewOAuth(secret, id, apiHost string) *oAuthApi {
	callbackUri := fmt.Sprintf("%s/user/login/google/callback", apiHost)
	scope := url.QueryEscape("https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/cloud-platform openid")
	authUri := fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&prompt=select_account", id, callbackUri, scope)

	return &oAuthApi{secret, id, authUri, callbackUri}
}

func (auth *oAuthApi) GetAuthUri() string {
	return auth.authUri
}

func (auth *oAuthApi) GetUserData(code string) (*oauth.AuthData, error) {
	respBody, err := auth.GetJSONResponse(code)
	if err != nil {
		return nil, err
	}

	var respMap map[string]interface{}

	if err := json.Unmarshal(respBody, &respMap); err != nil {
		return nil, err
	}

	email := respMap["email"].(string)
	idx := strings.Index(email, "@")

	return &oauth.AuthData{
		AccountID: respMap["id"].(string),
		Username:  email[:idx],
		Email:     email,
		Provider:  "Google",
	}, nil
}

func (auth *oAuthApi) GetJSONResponse(code string) ([]byte, error) {
	token, tokenType, err := auth.getAccessToken(code, "authorization_code")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"GET",
		"https://www.googleapis.com/oauth2/v1/userinfo?alt=json",
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", tokenType, token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if string(respBody) == "" {
		return nil, err
	}

	return respBody, nil
}

func (auth *oAuthApi) getAccessToken(codeOrToken, grantType string) (string, string, error) {
	queryName := "code"
	if grantType == "refresh_token" {
		queryName = grantType
	}

	reqBody, err := json.Marshal(map[string]string{
		"client_id":     auth.clientID,
		"client_secret": auth.clientSecret,
		queryName:       codeOrToken,
		"grant_type":    grantType,
		"redirect_uri":  auth.callbackUri,
	})
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequest(
		"POST",
		"https://oauth2.googleapis.com/token",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	if grantType == "refresh_token" {
		var idTokResp idTokenResponse
		if err := json.Unmarshal(respBody, &idTokResp); err != nil {
			return "", "", err
		}

		return idTokResp.AccessToken, idTokResp.TokenType, nil
	}

	var refTokResp refreshTokenResponse
	if err := json.Unmarshal(respBody, &refTokResp); err != nil {
		return "", "", err
	}

	token := refTokResp.AccessToken
	token_type := refTokResp.TokenType

	if refTokResp.Expires <= 0 {
		token, token_type, err = auth.getAccessToken(refTokResp.RefreshToken, "refresh_token")
		if err != nil {
			return "", "", err
		}
	}

	return token, token_type, nil
}
