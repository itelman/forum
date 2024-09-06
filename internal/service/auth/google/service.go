package google

import (
	"bytes"
	"encoding/json"
	"fmt"
	"forum/internal/repository/models"
	"forum/internal/service/auth"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func GetClientID() string {
	return os.Getenv("GOOGLE_CLIENT_ID")
}

func GetClientSecret() string {
	return os.Getenv("GOOGLE_CLIENT_SECRET")
}

func GetAccessToken(codeOrToken, grant_type string) (string, string, error) {
	clientID := GetClientID()
	clientSecret := GetClientSecret()

	queryName := "code"
	if grant_type == "refresh_token" {
		queryName = "refresh_token"
	}

	// Set us the request body as JSON
	requestBodyMap := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		queryName:       codeOrToken,
		"grant_type":    grant_type,
		"redirect_uri":  "https://forum-099y.onrender.com/auth/google/callback",
	}
	requestJSON, err := json.Marshal(requestBodyMap)
	if err != nil {
		return "", "", err
	}

	// POST request to set URL
	req, reqerr := http.NewRequest(
		"POST",
		"https://oauth2.googleapis.com/token",
		bytes.NewBuffer(requestJSON),
	)
	if reqerr != nil {
		return "", "", models.ErrBadGateway
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Get the response
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		return "", "", models.ErrBadGateway
	}

	// Response body converted to stringified JSON
	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	type AccessTokenResponse struct {
		AccessToken  string `json:"access_token"`
		Expires      int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
	}

	type NewAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		Expires     int    `json:"expires_in"`
		IdToken     string `json:"id_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	if grant_type == "refresh_token" {
		var ghresp NewAccessTokenResponse
		err = json.Unmarshal(respbody, &ghresp)
		if err != nil {
			return "", "", err
		}

		token := ghresp.AccessToken
		token_type := ghresp.TokenType

		return token, token_type, nil
	}

	// Convert stringified JSON to a struct object of type githubAccessTokenResponse
	var ghresp AccessTokenResponse
	err = json.Unmarshal(respbody, &ghresp)
	if err != nil {
		return "", "", err
	}

	token := ghresp.AccessToken
	token_type := ghresp.TokenType

	if ghresp.Expires <= 0 {
		token, token_type, err = GetAccessToken(ghresp.RefreshToken, "refresh_token")
		if err != nil {
			return "", "", err
		}
	}

	// Return the access token (as the rest of the
	// details are relatively unnecessary for us)
	return token, token_type, nil
}

func GetUserData(token_type, token string) (*auth.AuthData, error) {
	result := &auth.AuthData{}

	// Get request to a set URL
	req, reqerr := http.NewRequest(
		"GET",
		"https://www.googleapis.com/oauth2/v1/userinfo?alt=json",
		nil,
	)
	if reqerr != nil {
		return result, models.ErrBadGateway
	}

	// Set the Authorization header before sending the request
	// Authorization: token XXXXXXXXXXXXXXXXXXXXXXXXXXX
	authorizationHeaderValue := fmt.Sprintf("%s %s", token_type, token)
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

	result.AccountID = respMap["id"].(string)
	if err != nil {
		return result, err
	}

	result.Provider = "google_id"
	result.Email = respMap["email"].(string)

	idx := strings.Index(result.Email, "@")
	result.Username = result.Email[:idx]

	return result, nil
}
