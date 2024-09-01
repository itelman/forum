package sessions

import (
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
	UserID      int
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

func (s *SessionStore) PutSessionData(sessionID string, key string, value interface{}) {
	sessionData := s.GetSession(sessionID)
	if sessionData == nil {
		sessionData = &Session{Active: true}
	}

	if key == "flash" {
		sessionData.Flash = value.(string)
	} else if key == "userID" {
		sessionData.UserID = value.(int)
	}

	if key == "userID" {
		sessionData.LastRequest = time.Now()
	}

	s.PutSession(sessionID, sessionData)
}

func (s *SessionStore) CreateNewSession(userID ...int) (string, error) {
	s.SessionMutex.Lock()
	defer s.SessionMutex.Unlock()

	sessionID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	sessionIDStr := sessionID.String()

	if s.Store == nil {
		s.Store = make(map[string]*Session)
	}

	s.Store[sessionIDStr] = &Session{Active: true}
	if s.ActiveSessions == nil {
		s.ActiveSessions = make(map[int]string)
	}
	if len(userID) > 0 {
		s.ActiveSessions[userID[0]] = sessionIDStr
	}

	return sessionIDStr, nil
}

func (s *SessionStore) GetSession(sessionID string) *Session {
	sessionData, exists := s.Store[sessionID]
	if !exists {
		return &Session{Active: true}
	}
	return sessionData
}

func (s *SessionStore) PutSession(sessionID string, data *Session) {
	s.Store[sessionID] = data
}

func (s *SessionStore) DeleteSession(sessionID string) {
	s.SessionMutex.Lock()
	defer s.SessionMutex.Unlock()

	delete(s.Store, sessionID)
}

func (s *SessionStore) PopSessionFlash(sessionID string) string {
	sessionData, exists := s.Store[sessionID]
	if !exists {
		return ""
	}

	flash := sessionData.Flash
	sessionData.Flash = ""

	if sessionData.Flash == "" && sessionData.UserID <= 0 {
		delete(s.Store, sessionID)
	}

	return flash
}

func (s *SessionStore) GetSessionUserID(sessionID string) int {
	sessionData, exists := s.Store[sessionID]
	if !exists {
		return -1
	}

	userID := sessionData.UserID

	return userID
}

func (s *SessionStore) UpdateSessionLastReq(sessionID string) {
	s.Store[sessionID].LastRequest = time.Now()
}

func (s *SessionStore) DisableSession(sessionID string) {
	s.SessionMutex.Lock()
	defer s.SessionMutex.Unlock()

	s.Store[sessionID].Active = false
	delete(s.ActiveSessions, s.GetSessionUserID(sessionID))
}

func (s *SessionStore) GetSessionByUserID(userID int) (string, bool) {
	s.SessionMutex.Lock()
	defer s.SessionMutex.Unlock()

	sessionID, exists := s.ActiveSessions[userID]

	return sessionID, exists
}
