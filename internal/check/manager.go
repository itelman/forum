package check

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"forum/pkg/models"
	"forum/pkg/models/sqlite"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/golangcollege/sessions"
	_ "github.com/mattn/go-sqlite3"
)

type contextKey string

var contextKeyUser = contextKey("user")

type Config struct {
	Addr      string
	StaticDir string
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	session  *sessions.Session
	snippets interface {
		Insert(string, string, string, string, string) (int, error)
		Get(int) (*models.Snippet, error)
		Latest() ([]*models.Snippet, error)
		Filter(url.Values, func(int, string, string) (bool, error), func(int, []string, int) (bool, error)) ([]*models.Snippet, error)
		UpdateReactions(int, func(int) (int, error), func(int) (int, error)) error
		Paginate([]*models.Snippet, int, int) ([]*models.Snippet, int, error)
		// UpdateReactions(string, string, string, string, string) (int, int, error)
	}
	templatecache map[string]*template.Template
	users         interface {
		Insert(string, string, string, string) error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
	}
	comments interface {
		Insert(string, string, string) error
		Latest(int) ([]*models.Comment, error)
	}
	categories interface {
		Latest() ([]*models.Category, error)
	}
	post_category interface {
		Insert(string, []string) error
		Get(int) ([]string, error)
		FilterByCategories(int, []string, int) (bool, error)
	}
	post_reactions interface {
		Insert(string, string, string) error
		Get(string, string) (string, error)
		Delete(string, string) error
		Likes(int) (int, error)
		Dislikes(int) (int, error)
		FilterByLiked(int, string, string) (bool, error)
	}
	comment_reactions interface {
		Insert(string, string, string) error
		Get(string, string) (string, error)
		Delete(string, string) error
		Likes(int) (int, error)
		Dislikes(int) (int, error)
	}
}

func CreateServer(infoLog *log.Logger, errorLog *log.Logger, cfg *Config) *http.Server {
	dsn := flag.String("dsn", "db/forum.db?parseTime=true", "MySQL database")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	app := &application{
		errorLog:          errorLog,
		infoLog:           infoLog,
		session:           session,
		snippets:          &sqlite.SnippetModel{DB: db},
		templatecache:     templateCache,
		users:             &sqlite.UserModel{DB: db},
		comments:          &sqlite.CommentModel{DB: db},
		categories:        &sqlite.CategoryModel{DB: db},
		post_category:     &sqlite.PostCategoryModel{DB: db},
		post_reactions:    &sqlite.PostReactionModel{DB: db},
		comment_reactions: &sqlite.CommentReactionModel{DB: db},
	}
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
		},
	}
	srv := &http.Server{
		Addr:         cfg.Addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return srv
}

var db *sql.DB

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func CloseDB() error {
	return db.Close()
}
