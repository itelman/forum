package sessions

import (
	"forum/internal/service/auth"
	"net/http"
	"sync"
	"time"

	"github.com/gofrs/uuid"
)

type SessionStore struct {
	Store          map[string]*Session
	SessionMutex   sync.Mutex
	ActiveSessions map[int]string
}

type Session struct {
	Flash       string
	CSRFToken   string
	UserID      int
	AuthData    *auth.AuthData
	LastRequest time.Time
	Active      bool
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		Store:          make(map[string]*Session),
		ActiveSessions: make(map[int]string),
	}
}

func (s *SessionStore) GetSessionIDFromRequest(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			newSessionID, createErr := s.CreateNewSession()
			if createErr != nil {
				return "", createErr
			}
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    newSessionID,
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
			})
			return newSessionID, nil
		}
		return "", err
	}

	return cookie.Value, nil
}

func (s *SessionStore) PutSessionData(sessionID string, key string, value interface{}) {
	sessionData := s.GetSession(sessionID)
	if sessionData == nil {
		sessionData = &Session{Active: true}
	}

	if key == "flash" {
		sessionData.Flash = value.(string)
	} else if key == "authData" {
		sessionData.AuthData = value.(*auth.AuthData)
	} else if key == "csrf_token" {
		sessionData.CSRFToken = value.(string)
	}

	s.SessionMutex.Lock()
	s.Store[sessionID] = sessionData
	s.SessionMutex.Unlock()
}

func (s *SessionStore) CreateNewSession(userID ...int) (string, error) {
	sessionID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	sessionIDStr := sessionID.String()

	s.SessionMutex.Lock()

	s.Store[sessionIDStr] = &Session{Active: true}

	if len(userID) > 0 {
		s.Store[sessionIDStr].UserID = userID[0]
		s.Store[sessionIDStr].LastRequest = time.Now()
		s.ActiveSessions[userID[0]] = sessionIDStr
	}

	s.SessionMutex.Unlock()

	return sessionIDStr, nil
}

func (s *SessionStore) NewSessionID() (string, error) {
	sessionID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return sessionID.String(), nil
}

func (s *SessionStore) GetSession(sessionID string) *Session {
	s.SessionMutex.Lock()
	sessionData, exists := s.Store[sessionID]
	s.SessionMutex.Unlock()

	if !exists {
		return nil
	}

	return sessionData
}

func (s *SessionStore) DeleteSession(sessionID string) {
	s.SessionMutex.Lock()
	defer s.SessionMutex.Unlock()

	if _, exists := s.Store[sessionID]; exists {
		delete(s.ActiveSessions, s.Store[sessionID].UserID)
	}

	delete(s.Store, sessionID)
}

func (s *SessionStore) PopSessionFlash(sessionID string) string {
	s.SessionMutex.Lock()
	defer s.SessionMutex.Unlock()

	sessionData, exists := s.Store[sessionID]
	if !exists {
		return ""
	}

	flash := sessionData.Flash
	s.Store[sessionID].Flash = ""

	return flash
}

func (s *SessionStore) UpdateSessionLastReq(sessionID string) {
	s.SessionMutex.Lock()
	s.Store[sessionID].LastRequest = time.Now()
	s.SessionMutex.Unlock()
}

func (s *SessionStore) DisableSession(sessionID string) {
	s.SessionMutex.Lock()
	s.Store[sessionID].Active = false

	delete(s.ActiveSessions, s.Store[sessionID].UserID)
	s.SessionMutex.Unlock()
}

func (s *SessionStore) GetSessionByUserID(userID int) (string, bool) {
	s.SessionMutex.Lock()
	defer s.SessionMutex.Unlock()

	sessionID, exists := s.ActiveSessions[userID]

	return sessionID, exists
}
