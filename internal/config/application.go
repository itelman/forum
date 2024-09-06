package config

import (
	"database/sql"
	"forum/internal/repository"
	"forum/internal/repository/sqlite"
	"forum/pkg/sessions"
	"forum/pkg/tmplcache"
	"html/template"
	"log"
	"time"
)

type Application struct {
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	CookieLimit   time.Duration
	SessionStore  *sessions.SessionStore
	TemplateCache tmplcache.TemplateCache
	Repository    *repository.Repository
}

func NewApplication(infoLog *log.Logger, errorLog *log.Logger, db *sql.DB, templateCache map[string]*template.Template) *Application {
	return &Application{
		InfoLog:       infoLog,
		ErrorLog:      errorLog,
		CookieLimit:   15 * time.Minute,
		SessionStore:  sessions.NewSessionStore(),
		TemplateCache: templateCache,
		Repository:    sqlite.NewRepository(db),
	}
}
