package main

import (
	"forum/internal/config"
	"forum/internal/handler"
	"forum/internal/router"
	"forum/pkg/store"
	"forum/pkg/tmplcache"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer f.Close()

	db, err := store.NewSqlite("./storage/storage.db?parseTime=true")
	if err != nil {
		errorLog.Fatal(err)
	}

	err = store.Migrate(db, "./migrations/sqlite/00001_initial.up.sql")
	if err != nil {
		errorLog.Fatal(err)
	}

	templateCache, err := tmplcache.NewTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := config.NewApplication(infoLog, errorLog, db, templateCache)

	handlers := &handler.Handlers{App: app}

	srv := &http.Server{
		Addr:         ":8080",
		ErrorLog:     errorLog,
		Handler:      router.Router(handlers),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on http://localhost%s", srv.Addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
