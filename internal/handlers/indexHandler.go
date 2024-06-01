package handlers

import (
	"forum/internal/models"
	"net/http"
	"text/template"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "IndexHandler"

	if r.URL.Path != "/" {
		ErrorHandler(w, http.StatusNotFound, handlerName, "Endpoint Failed")
		return
	}

	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed, handlerName, "GET/POST Method Failed")
		return
	}

	db, err := models.StorageConstructor("internal/storage/storage.db")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}

	posts, err := db.GetAllPosts()
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}

	/* this block of code needs to be internalized somehow */
	for _, post := range posts {
		post.SetReaction(&db, 1)
	}
	/* */

	results, err := db.SetResults(posts)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}

	tmpl, err := template.ParseFiles("templates/index/SAMPLE.html", "templates/index/index.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}

	err = tmpl.Execute(w, results)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}
}

func ResultsHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "ResultsHandler"

	if r.URL.Path != "/results" {
		ErrorHandler(w, http.StatusNotFound, handlerName, "Endpoint Failed")
		return
	}

	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed, handlerName, "GET/POST Method Failed")
		return
	}

	db, err := models.StorageConstructor("internal/storage/storage.db")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}

	r.ParseForm()
	posts, err := models.GetPostsResults(r.Form, &db)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}

	results, err := db.SetResults(posts)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}

	tmpl, err := template.ParseFiles("templates/index/SAMPLE.html", "templates/index/index.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}

	err = tmpl.Execute(w, results)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}
}
