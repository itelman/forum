package handler

import (
	"forum/internal/config"
	"forum/pkg/tmplcache"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Handlers struct {
	App *config.Application
}

func MockHandlers() *Handlers {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	path, err := os.Getwd()
	if err != nil {
		errorLog.Fatal(err)
	}
	workDir := strings.ReplaceAll(path, "/internal/handler", "")
	templateCache, err := tmplcache.NewTemplateCache(filepath.Join(workDir, "/ui/html/"))
	if err != nil {
		errorLog.Fatal(err)
	}
	app := config.MockApplication(errorLog, templateCache)
	return &Handlers{app}
}
