package config

import (
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

func (app *Application) GetSessionIDFromRequest(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			newSessionID, createErr := app.CreateNewSession()
			if createErr != nil {
				return "", createErr
			}
			http.SetCookie(w, &http.Cookie{
				Name:  "session_id",
				Value: newSessionID,
				Path:  "/",
			})
			return newSessionID, nil
		}
		return "", err
	}

	return cookie.Value, nil
}

func (app *Application) PutSessionData(sessionID string, key string, value interface{}) {
	sessionData := app.GetSession(sessionID)
	if sessionData == nil {
		sessionData = make(map[string]interface{})
	}
	sessionData[key] = value
	if key == "userID" {
		sessionData["lastRequest"] = time.Now()
	}
	app.PutSession(sessionID, sessionData)
}

func (app *Application) CreateNewSession(userID ...int) (string, error) {
	sessionID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	sessionIDStr := sessionID.String()

	app.SessionMutex.Lock()
	defer app.SessionMutex.Unlock()

	if app.SessionStore == nil {
		app.SessionStore = make(map[string]map[string]interface{})
	}

	if len(userID) > 0 {
		existingSessionID, exists := app.ActiveSessions[userID[0]]
		if exists {
			app.DeleteSession(existingSessionID)
		}
	}

	app.SessionStore[sessionIDStr] = make(map[string]interface{})
	if app.ActiveSessions == nil {
		app.ActiveSessions = make(map[int]string)
	}
	if len(userID) > 0 {
		app.ActiveSessions[userID[0]] = sessionIDStr
	}

	return sessionIDStr, nil
}

func (app *Application) GetSession(sessionID string) map[string]interface{} {
	app.SessionMutex.Lock()
	defer app.SessionMutex.Unlock()
	sessionData, exists := app.SessionStore[sessionID]
	if !exists {
		return make(map[string]interface{})
	}
	return sessionData
}

func (app *Application) PutSession(sessionID string, data map[string]interface{}) {
	app.SessionMutex.Lock()
	defer app.SessionMutex.Unlock()
	app.SessionStore[sessionID] = data
}

func (app *Application) DeleteSession(sessionID string) {
	app.SessionMutex.Lock()
	defer app.SessionMutex.Unlock()
	delete(app.SessionStore, sessionID)
}

func (app *Application) SetSessionToken(sessionID, token string) {
	app.SessionMutex.Lock()
	defer app.SessionMutex.Unlock()

	session, exists := app.SessionStore[sessionID]
	if !exists {
		session = make(map[string]interface{})
		app.SessionStore[sessionID] = session
	}
	session["csrf_token"] = token
}

func (app *Application) GetSessionToken(sessionID string) (string, bool) {
	app.SessionMutex.Lock()
	defer app.SessionMutex.Unlock()

	sessionData, exists := app.SessionStore[sessionID]
	if !exists {
		return "", false
	}

	token, exists := sessionData["csrf_token"]
	if !exists {
		return "", false
	}

	csrfToken, ok := token.(string)
	if !ok {
		return "", false
	}

	return csrfToken, true
}

func (app *Application) DeleteSessionData(sessionID string, key string) {
	app.SessionMutex.Lock()
	defer app.SessionMutex.Unlock()

	sessionData, exists := app.SessionStore[sessionID]
	if !exists {
		return
	}
	delete(sessionData, key)
	if len(sessionData) == 0 {
		delete(app.SessionStore, sessionID)
	}
}

func (app *Application) GetSessionUserID(sessionID string) int {
	app.SessionMutex.Lock()
	defer app.SessionMutex.Unlock()

	sessionData, exists := app.SessionStore[sessionID]
	if !exists {
		return 0
	}

	userID, ok := sessionData["userID"].(int)
	if !ok {
		return 0
	}

	return userID
}

func (app *Application) UpdateSessionLastReq(sessionID string) {
	app.SessionStore[sessionID]["lastRequest"] = time.Now()
}
