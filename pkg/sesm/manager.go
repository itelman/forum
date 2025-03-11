package sesm

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"
)

const (
	SessionId      = "session_id"
	UserId         = "user_id"
	LastRequest    = "last_request"
	Status         = "status"
	BlockTimestamp = "block_timestamp"
	Flash          = "flash"
)

var (
	ErrSessionNotFound = errors.New("SESM: session not found")
	ErrDataNotFound    = errors.New("SESM: session data not found")
)

type session map[string]interface{}

type SessionManager interface {
	CreateSession(userId int) (string, error)
	CurrentSessionID(r *http.Request) (string, error)
	DeleteCurrentSession(r *http.Request) error
	GetSessionData(r *http.Request, key string) (interface{}, error)
	AddOrUpdateSessionData(r *http.Request, dataMap session) error
	DeleteSessionData(r *http.Request, keys []string) error
	UpdateSessionLastRequest(r *http.Request) error
	UpdateSessionFlash(r *http.Request, val string) error
	PopSessionFlash(r *http.Request) (string, error)
	DeleteAllUserSessions(userId int)
}

type sessionManager struct {
	store       map[string]session
	activeUsers map[int]map[string]interface{}
	mutex       sync.RWMutex
}

func NewSessionManager() *sessionManager {
	return &sessionManager{
		store:       make(map[string]session),
		activeUsers: make(map[int]map[string]interface{}),
	}
}

func (s *sessionManager) CreateSession(userId int) (string, error) {
	newUUID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	sessionId := newUUID.String()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.store[sessionId] = session{UserId: userId, LastRequest: time.Now(), Status: "active"}

	if len(s.activeUsers[userId]) == 0 {
		s.activeUsers[userId] = make(map[string]interface{})
	}
	s.activeUsers[userId][sessionId] = nil

	return sessionId, nil
}

func (s *sessionManager) CurrentSessionID(r *http.Request) (string, error) {
	cookie, err := r.Cookie(SessionId)
	if err != nil {
		return "", ErrSessionNotFound
	}

	s.mutex.RLock() // Read lock allows multiple readers
	defer s.mutex.RUnlock()

	_, exists := s.store[cookie.Value]
	if !exists {
		return "", ErrSessionNotFound
	}

	return cookie.Value, nil
}

func (s *sessionManager) DeleteCurrentSession(r *http.Request) error {
	sessionID, err := s.CurrentSessionID(r)
	if err != nil {
		return nil
	}

	s.mutex.RLock()
	uid, exists := s.store[sessionID][UserId]
	if !exists {
		return ErrDataNotFound
	}
	s.mutex.RUnlock()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.store, sessionID)
	delete(s.activeUsers[uid.(int)], sessionID)

	return nil
}

func (s *sessionManager) GetSessionData(r *http.Request, key string) (interface{}, error) {
	sessionId, err := s.CurrentSessionID(r)
	if err != nil {
		return nil, err
	}

	s.mutex.RLock() // Read lock allows multiple readers
	defer s.mutex.RUnlock()

	val, exists := s.store[sessionId][key]
	if !exists {
		return nil, ErrDataNotFound
	}

	return val, nil
}

func (s *sessionManager) AddOrUpdateSessionData(r *http.Request, dataMap session) error {
	sessionID, err := s.CurrentSessionID(r)
	if err != nil {
		return err
	}

	s.mutex.Lock() // Locking for exclusive write access
	defer s.mutex.Unlock()

	for key, val := range dataMap {
		s.store[sessionID][key] = val
	}

	return nil
}

func (s *sessionManager) DeleteSessionData(r *http.Request, keys []string) error {
	sessionID, err := s.CurrentSessionID(r)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, key := range keys {
		delete(s.store[sessionID], key)
	}

	return nil
}

func (s *sessionManager) UpdateSessionLastRequest(r *http.Request) error {
	if err := s.AddOrUpdateSessionData(r, session{LastRequest: time.Now()}); err != nil {
		return err
	}

	return nil
}

func (s *sessionManager) UpdateSessionFlash(r *http.Request, val string) error {
	if err := s.AddOrUpdateSessionData(r, session{Flash: val}); err != nil {
		return err
	}

	return nil
}

func (s *sessionManager) PopSessionFlash(r *http.Request) (string, error) {
	flash, err := s.GetSessionData(r, Flash)
	if errors.Is(err, ErrDataNotFound) {
		return "", nil
	} else if err != nil {
		return "", err
	}

	if err := s.DeleteSessionData(r, []string{Flash}); err != nil {
		return "", err
	}

	return flash.(string), nil
}

func (s *sessionManager) DeleteAllUserSessions(userId int) {
	for sessionId, _ := range s.activeUsers[userId] {
		delete(s.store, sessionId)
	}

	s.activeUsers[userId] = make(map[string]interface{})
}
