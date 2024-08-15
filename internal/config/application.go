package config

import (
	"database/sql"
	"forum/internal/repository/models"
	"html/template"
	"io/ioutil"
	"log"
	"net/url"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type contextKey string

var contextKeyUser = contextKey("user")

type Config struct {
	Addr      string
	StaticDir string
}

type Application struct {
	ErrorLog       *log.Logger
	InfoLog        *log.Logger
	SessionStore   map[string]map[string]interface{}
	SessionMutex   sync.Mutex
	ActiveSessions map[int]string
	Templatecache  map[string]*template.Template
	Posts          interface {
		Insert(string, string, string) (int, error)
		Get(int) (*models.Post, error)
		Latest() ([]*models.Post, error)
		Filter(url.Values, func(int, string, string) (bool, error), func(int, []string, int) (bool, error)) ([]*models.Post, error)
		UpdateReactions(int, func(int) (int, error), func(int) (int, error)) error
	}
	Users interface {
		Insert(string, string, string) error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
	}
	Comments interface {
		Insert(string, string, string) error
		Latest(int) ([]*models.Comment, error)
		UpdateReactions(int, func(int) (int, error), func(int) (int, error)) error
	}
	Categories interface {
		Latest() ([]*models.Category, error)
	}
	Post_category interface {
		Insert(string, []string) error
		Get(int) ([]string, error)
		FilterByCategories(int, []string, int) (bool, error)
	}
	Post_reactions interface {
		Insert(string, string, string) error
		Get(string, string) (string, error)
		Delete(string, string) error
		Likes(int) (int, error)
		Dislikes(int) (int, error)
		FilterByLiked(int, string, string) (bool, error)
	}
	Comment_reactions interface {
		Insert(string, string, string) error
		Get(string, string) (string, error)
		Delete(string, string) error
		Likes(int) (int, error)
		Dislikes(int) (int, error)
	}
}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	filePath := "./migrations/sqlite/00001_initial.up.sql"
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	script := string(content)

	_, err = db.Exec(script)
	if err != nil {
		return nil, err
	}

	return db, nil
}
