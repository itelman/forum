package router

import (
	"forum/internal/handler"
	"forum/internal/service/middleware"
	"net/http"
)

type routes struct {
	Path        string
	Handler     func(http.ResponseWriter, *http.Request)
	RequireAuth bool
}

func Router(handlers *handler.Handlers) http.Handler {
	mux := http.NewServeMux()

	routes := []routes{
		{"/", handlers.Home, false},
		{"/results", handlers.Results, false},
		{"/post", handlers.ShowPost, false},
		{"/post/create", handlers.CreatePostForm, true},
		{"/post/comment", handlers.CreateComment, true},
		{"/post/reaction", handlers.HandlePostReaction, true},
		{"/comment/reaction", handlers.HandleCommentReaction, true},
		{"/user/signup", handlers.SignupUserForm, false},
		{"/user/login", handlers.LoginUserForm, false},
		{"/user/logout", handlers.LogoutUser, true},
	}

	middleware := &middleware.Middleware{Handlers: handlers}

	for _, route := range routes {
		mux.Handle(route.Path, middleware.DynamicMiddleware(route.Handler, route.RequireAuth))
	}

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	return middleware.StandardMiddleware(mux)
}