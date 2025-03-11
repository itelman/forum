package templates

import (
	"html/template"
	"path/filepath"
	"strings"
	"time"
)

type TemplateCache map[string]*template.Template

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Local().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func NewTemplateCache(dir string) (TemplateCache, error) {
	cache := TemplateCache{}

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

		name = strings.ReplaceAll(name, ".html", "")
		cache[name] = ts
	}

	return cache, nil
}
