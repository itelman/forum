package dto

import "net/http"

type methods []string

var (
	GetMethod      = methods{http.MethodGet}
	PostMethod     = methods{http.MethodPost}
	GetPostMethods = methods{http.MethodGet, http.MethodPost}
)

type Route struct {
	Path    string
	Methods methods
	Handler func(http.ResponseWriter, *http.Request)
}
