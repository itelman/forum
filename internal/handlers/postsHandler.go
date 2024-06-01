package handlers

import (
	"database/sql"
	"forum/internal/models"
	"net/http"
	"strconv"
	"text/template"
)

var token string

func NewPostHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "NewPostHandler"

	if r.URL.Path != "/posts/new" {
		ErrorHandler(w, http.StatusNotFound, handlerName, "Endpoint Failed")
		return
	}

	if r.Method == http.MethodPost {
		PostVerifyHandler(w, r, 1)
		return
	}

	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed, handlerName, "GET/POST Method Failed")
		return
	}

	tmpl, err := template.ParseFiles("templates/index/SAMPLE.html", "templates/index/newPost.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
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

	categories, err := db.GetAllCategories()
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}

	err = tmpl.Execute(w, categories)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}
}

func PostVerifyHandler(w http.ResponseWriter, r *http.Request, user_id int) {
	handlerName := "PostVerifyHandler"

	post := models.Post{}
	var categories_id []string

	r.ParseForm()
	for k, v := range r.Form {
		switch k {
		case "title":
			post.Title = v[0]
		case "content":
			post.Content = v[0]
		case "category":
			categories_id = v
		}
	}
	post.UserID = user_id

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

	id, err := db.SetPostID()
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}
	post.ID = id

	err = db.InsertPost(post, categories_id)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}

	tmpl, err := template.ParseFiles("templates/index/SAMPLE.html", "templates/index/postOK.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}

	err = tmpl.Execute(w, post.ID)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}
}

func ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "ViewPostHandler"

	if r.URL.Path != "/posts" {
		ErrorHandler(w, http.StatusNotFound, handlerName, "Endpoint Failed")
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

	if r.Method == http.MethodPost {
		r.ParseForm()

		if _, ok := r.Form["comment"]; ok {
			CommentVerifyHandler(w, r, &db, 1)
		} else if _, ok := r.Form["postReaction"]; ok {
			PostReactionHandler(w, r, &db, 1)
		} else if _, ok := r.Form["commentReaction"]; ok {
			CommentReactionHandler(w, r, &db, 1)
		} else {
			ErrorHandler(w, http.StatusBadRequest, handlerName, err.Error())
			return
		}
	}

	if !(r.Method == http.MethodGet || r.Method == http.MethodPost) {
		ErrorHandler(w, http.StatusMethodNotAllowed, handlerName, "GET/POST Method Failed")
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, handlerName, err.Error())
		return
	}

	post, err := db.GetPostByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			ErrorHandler(w, http.StatusNotFound, handlerName, err.Error())
		} else {
			ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		}
		return
	}

	/* this block of code needs to be internalized somehow */
	post.SetReaction(&db, 1)

	for _, comment := range post.Comments {
		comment.SetReaction(&db, 1)
	}
	/* */

	tmpl, err := template.ParseFiles("templates/index/SAMPLE.html", "templates/index/viewPost.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}

	err = tmpl.Execute(w, post)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}
}
