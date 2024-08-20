package handler

import (
	"forum/internal/config"
	"forum/internal/repository/mock"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Handlers struct {
	App *config.Application
}

func MockHandlers() *Handlers {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	path, err := os.Getwd()
	if err != nil {
		errorLog.Fatal(err)
	}

	workDir := strings.ReplaceAll(path, "/internal/handler", "")

	templateCache, err := config.NewTemplateCache(filepath.Join(workDir, "/ui/html/"))
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &config.Application{
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
		CookieLimit:   15 * time.Minute,
		SessionStore:  make(map[string]map[string]interface{}),
		Posts:         mock.NewPostModel(),
		Templatecache: templateCache,
		Comments:      mock.NewCommentModel(),
		Post_category: mock.NewPostCategoryModel(),
	}

	handlers := &Handlers{app}

	return handlers
}
