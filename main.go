package main

import (
	"forum/internal/app"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	srv := app.CreateServer(infoLog, errorLog)

	infoLog.Printf("Starting server on http://localhost%s", ":8080")
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
