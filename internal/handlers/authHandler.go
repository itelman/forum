package handlers

import (
	"forum/internal/models"
	"net/http"
	"text/template"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "AuthHandler"

	if r.URL.Path != "/auth" {
		ErrorHandler(w, http.StatusNotFound, handlerName, "endpoint")
		return
	}

	if r.Method == http.MethodPost {
		AuthVerifyHandler(w, r)
		return
	}

	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed, handlerName, "GET/POST Method Failed")
		return
	}

	AuthTmplExecute(w, handlerName, "")
}

func AuthTmplExecute(w http.ResponseWriter, handlerName, msg string) {
	tmpl, err := template.ParseFiles("templates/sign-in/SAMPLE.html", "templates/sign-in/auth.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, "tmplParse")
		return
	}

	err = tmpl.Execute(w, msg)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, "tmplExecute")
		return
	}
}

func AuthVerifyHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "AuthVerifyHandler"

	user := models.User{}

	r.ParseForm()
	for k, v := range r.Form {
		switch k {
		case "username":
			user.Username = v[0]
		case "password":
			user.Password = v[0]
		}
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

	ok, err := db.CheckCredentials(user)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}

	if !ok {
		msg := "(!) Incorrect username and/or password. Please make sure you entered correct credentials and try again."
		AuthTmplExecute(w, handlerName, msg)
	}

	tmpl, err := template.ParseFiles("templates/index/SAMPLE.html", "templates/index/authOK.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		return
	}
}
