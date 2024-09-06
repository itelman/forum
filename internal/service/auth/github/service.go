package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"forum/internal/repository/models"
	"forum/internal/service/auth"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func GetClientID() string {
	return os.Getenv("GITHUB_CLIENT_ID")
}

func GetClientSecret() string {
	return os.Getenv("GITHUB_CLIENT_SECRET")
}

func GetAccessToken(code string) (string, error) {
	clientID := GetClientID()
	clientSecret := GetClientSecret()

	// Set us the request body as JSON
	requestBodyMap := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
	}
	requestJSON, err := json.Marshal(requestBodyMap)
	if err != nil {
		return "", err
	}

	// POST request to set URL
	req, reqerr := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON),
	)
	if reqerr != nil {
		return "", models.ErrBadGateway
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Get the response
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		return "", models.ErrBadGateway
	}

	// Response body converted to stringified JSON
	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Represents the response received from Github
	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	// Convert stringified JSON to a struct object of type githubAccessTokenResponse
	var ghresp githubAccessTokenResponse
	err = json.Unmarshal(respbody, &ghresp)
	if err != nil {
		return "", err
	}

	// Return the access token (as the rest of the
	// details are relatively unnecessary for us)
	return ghresp.AccessToken, nil
}

func GetUserData(accessToken string) (*auth.AuthData, error) {
	result := &auth.AuthData{}

	// Get request to a set URL
	req, reqerr := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if reqerr != nil {
		return result, models.ErrBadGateway
	}

	// Set the Authorization header before sending the request
	// Authorization: token XXXXXXXXXXXXXXXXXXXXXXXXXXX
	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	// Make the request
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		return result, models.ErrBadGateway
	}

	// Read the response as a byte slice
	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	if string(respbody) == "" {
		return result, models.ErrBadGateway
	}

	var respMap map[string]interface{}

	// Unmarshal the JSON into the map
	err = json.Unmarshal(respbody, &respMap)
	if err != nil {
		return result, err
	}

	result.AccountID = strconv.Itoa(int(respMap["id"].(float64)))
	result.Provider = "github_id"

	if val, ok := respMap["login"]; !ok || val == nil {
		result.Username = ""
	} else {
		result.Username = respMap["login"].(string)
	}

	if val, ok := respMap["email"]; !ok || val == nil {
		result.Email = ""
	} else {
		result.Email = respMap["email"].(string)
	}

	return result, nil
}
