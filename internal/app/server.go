package app

import (
	"flag"
	"forum/internal/config"
	"forum/internal/handler"
	"forum/internal/repository/sqlite"
	"log"
	"net/http"
	"os"
	"time"
)

func CreateServer(infoLog *log.Logger, errorLog *log.Logger) *http.Server {
	dsn := flag.String("dsn", "pkg/store/forum.db?parseTime=true", "MySQL database")
	f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	db, err := config.OpenDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	templateCache, err := config.NewTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &config.Application{
		ErrorLog:          errorLog,
		InfoLog:           infoLog,
		CookieLimit:       15 * time.Minute,
		SessionStore:      make(map[string]map[string]interface{}),
		Posts:             &sqlite.PostModel{DB: db},
		Templatecache:     templateCache,
		Users:             &sqlite.UserModel{DB: db},
		Comments:          &sqlite.CommentModel{DB: db},
		Categories:        &sqlite.CategoryModel{DB: db},
		Post_category:     &sqlite.PostCategoryModel{DB: db},
		Post_reactions:    &sqlite.PostReactionModel{DB: db},
		Comment_reactions: &sqlite.CommentReactionModel{DB: db},
	}

	handlers := &handler.Handlers{App: app}

	srv := &http.Server{
		Addr:         ":8080",
		ErrorLog:     errorLog,
		Handler:      Router(handlers),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return srv
}
