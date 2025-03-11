package main

import (
	"database/sql"
	"github.com/itelman/forum/pkg/oauth"
	"github.com/itelman/forum/pkg/oauth/github"
	"github.com/itelman/forum/pkg/oauth/google"
	"github.com/itelman/forum/pkg/sesm"
	"github.com/itelman/forum/pkg/sqlite"
	"github.com/itelman/forum/pkg/templates"
)

type Dependencies struct {
	sqlite        *sql.DB
	githubAuth    oauth.AuthApi
	googleAuth    oauth.AuthApi
	sesManager    sesm.SessionManager
	templateCache templates.TemplateCache
}

func (d *Dependencies) Close() {
	if d.sqlite != nil {
		d.sqlite.Close()
	}
}

func NewDependencies(opts ...Option) (deps *Dependencies, err error) {
	deps = &Dependencies{}
	for _, opt := range opts {
		if err := opt(deps); err != nil {
			return nil, err
		}
	}

	deps.sesManager = sesm.NewSessionManager()

	return deps, nil
}

type Option func(*Dependencies) error

func WithSqlite(dbDir, migrDir string) Option {
	return func(d *Dependencies) error {
		db, err := sqlite.NewSqlite(dbDir)
		if err != nil {
			return err
		}

		if err := sqlite.Migrate(db, migrDir); err != nil {
			return err
		}

		d.sqlite = db
		return nil
	}
}

func WithGithubAuth(clientSecret, clientId, apiHost string) Option {
	return func(d *Dependencies) error {
		d.githubAuth = github.NewOAuth(clientSecret, clientId, apiHost)
		return nil
	}
}

func WithGoogleAuth(clientSecret, clientId, apiHost string) Option {
	return func(d *Dependencies) error {
		d.googleAuth = google.NewOAuth(clientSecret, clientId, apiHost)
		return nil
	}
}

func WithTemplateCache(dir string) Option {
	return func(d *Dependencies) error {
		templateCache, err := templates.NewTemplateCache(dir)
		if err != nil {
			return err
		}

		d.templateCache = templateCache
		return nil
	}
}
