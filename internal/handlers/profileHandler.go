package handlers

import (
	"forum/internal/models"
	"net/http"
	"strconv"
	"text/template"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "ProfileHandler"

	if r.URL.Path != "/profile" {
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

	if !(r.Method == http.MethodGet) {
		ErrorHandler(w, http.StatusMethodNotAllowed, handlerName, "GET/POST Method Failed")
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, handlerName, err.Error())
		return
	}

	user, err := db.GetUserBy("id", id)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}

	tmpl, err := template.ParseFiles("templates/index/SAMPLE.html", "templates/index/profile.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}

	err = tmpl.Execute(w, user)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}
}
