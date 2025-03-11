package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/itelman/forum/pkg/oauth"
	"io/ioutil"
	"net/http"
	"strconv"
)

type githubAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type oAuthApi struct {
	clientSecret string
	clientID     string
	authUri      string
}

func NewOAuth(secret, id, apiHost string) *oAuthApi {
	authUri := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s/user/login/github/callback", id, apiHost)
	return &oAuthApi{secret, id, authUri}
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

	var username string
	if val, ok := respMap["login"]; !ok || val == nil {
		username = ""
	} else {
		username = respMap["login"].(string)
	}

	var email string
	if val, ok := respMap["email"]; !ok || val == nil {
		email = ""
	} else {
		email = respMap["email"].(string)
	}

	return &oauth.AuthData{
		AccountID: strconv.Itoa(int(respMap["id"].(float64))),
		Username:  username,
		Email:     email,
		Provider:  "Github",
	}, nil
}

func (auth *oAuthApi) GetJSONResponse(code string) ([]byte, error) {
	token, err := auth.getAccessToken(code)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))

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

func (auth *oAuthApi) getAccessToken(code string) (string, error) {
	reqBody, err := json.Marshal(map[string]string{
		"client_id":     auth.clientID,
		"client_secret": auth.clientSecret,
		"code":          code,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var githubResp githubAccessTokenResponse
	if err := json.Unmarshal(respBody, &githubResp); err != nil {
		return "", err
	}

	return githubResp.AccessToken, nil
}
