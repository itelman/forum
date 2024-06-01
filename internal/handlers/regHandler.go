package handlers

import (
	"forum/internal/models"
	"net/http"
	"text/template"
)

func RegHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "RegHandler"

	if r.URL.Path != "/signup" {
		ErrorHandler(w, http.StatusNotFound, handlerName, "endpoint")
		return
	}

	if r.Method == http.MethodPost {
		RegVerifyHandler(w, r)
		return
	}

	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed, handlerName, "GET/POST Method Failed")
		return
	}

	RegTmplExecute(w, handlerName, "")
}

func RegTmplExecute(w http.ResponseWriter, handlerName, msg string) {
	tmpl, err := template.ParseFiles("templates/sign-in/SAMPLE.html", "templates/sign-in/signup.html")
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

func RegVerifyHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "RegVerifyHandler"

	user := models.User{}

	r.ParseForm()
	for k, v := range r.Form {
		switch k {
		case "email":
			user.Email = v[0]
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

	err = db.InsertUser(user)
	if err != nil {
		if err.Error() == "Email/username exists" {
			msg := "(!) The provided email and/or username is linked to another account. Please try again with a different email or username."
			RegTmplExecute(w, handlerName, msg)
		} else {
			ErrorHandler(w, http.StatusInternalServerError, handlerName, err.Error())
		}
		return
	}

	tmpl, err := template.ParseFiles("templates/index/SAMPLE.html", "templates/index/signupOK.html")
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
