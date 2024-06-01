package main

import (
	"fmt"
	"forum/internal/handlers"
	"log"
	"net/http"
)

func GetVal(value interface{}) {
	fmt.Printf("%#v\n", value)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.IndexHandler)
	mux.HandleFunc("/results", handlers.ResultsHandler)
	mux.HandleFunc("/signup", handlers.RegHandler)
	mux.HandleFunc("/auth", handlers.AuthHandler)
	mux.HandleFunc("/profile", handlers.ProfileHandler)
	mux.HandleFunc("/posts/new", handlers.NewPostHandler)
	mux.HandleFunc("/posts", handlers.ViewPostHandler)
	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))

	fmt.Println("Server is running on http://localhost:8000/...")
	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
