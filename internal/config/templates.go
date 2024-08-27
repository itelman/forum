package config

import (
	"forum/internal/repository/models"
	"forum/internal/service/forms"
	"html/template"
	"path/filepath"
	"time"
)

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

type TemplateData struct {
	AuthenticatedUser *models.User
	CurrentYear       int
	Flash             string
	Form              *forms.Form
	Post              *models.Post
	Posts             []*models.Post
	Comments          []*models.Comment
	Categories        []*models.Category
	PCRelations       []string
	Error             *models.Error
}

func NewTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(filepath.Join(dir, "*page.html"))
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob(filepath.Join(dir, "*layout.html"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*partial.html"))
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
