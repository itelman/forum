package handlers

import (
	"fmt"
	"forum/internal/models"
	"log"
	"net/http"
	"text/template"
)

func ErrorHandler(w http.ResponseWriter, statusCode int, handlerName, msg string) {
	handlerErr := "ErrorHandler"

	tmpl, err := template.ParseFiles("templates/index/SAMPLE.html", "templates/index/error.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(GetHandlerMsg(models.ErrorConstructor(http.StatusInternalServerError), handlerErr, err.Error()))
		return
	}

	w.WriteHeader(statusCode)
	errStr := models.ErrorConstructor(statusCode)
	err = tmpl.Execute(w, errStr)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(GetHandlerMsg(models.ErrorConstructor(http.StatusInternalServerError), handlerErr, err.Error()))
		return
	}
	log.Println(GetHandlerMsg(errStr, handlerName, msg))
}

func GetHandlerMsg(errStr models.Error, handlerName, msg string) string {
	return fmt.Sprintf("\"[CODE %d]: %s\" - %s: %s", errStr.Code, errStr.Name, handlerName, msg)
}
