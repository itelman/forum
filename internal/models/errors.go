package models

import "net/http"

type Error struct {
	Code    int
	Name    string
	Content string
}

func ErrorConstructor(statusCode int) Error {
	errContent := map[int]string{
		http.StatusBadRequest:          "The server cannot process the request due to something that is perceived to be a client error.",
		http.StatusNotFound:            "The requested resource could not be found but may be available again in the future.",
		http.StatusMethodNotAllowed:    "A request method is not supported for the requested resource.",
		http.StatusInternalServerError: "An unexpected condition was encountered.<br />Our service team has been dispatched to bring it back online.",
	}

	return Error{Code: statusCode, Name: http.StatusText(statusCode), Content: errContent[statusCode]}
}
