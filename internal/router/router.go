package router

import (
	"forum/internal/handler"
	"forum/internal/service/middleware"
	"net/http"
	"time"
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
		{"/post/delete", handlers.DeletePost, true},
		{"/post/edit", handlers.EditPostForm, true},
		{"/post/comment/delete", handlers.DeleteComment, true},
		{"/post/comment/edit", handlers.EditCommentForm, true},
		{"/user/activity/created", handlers.ShowCreatedPosts, true},
		{"/user/activity/reacted", handlers.ShowReactedPosts, true},
		{"/user/activity/commented", handlers.ShowCommentedPosts, true},
		{"/user/notifications", handlers.ShowNotifications, true},
		{"/user/signup/provider", handlers.SignupUserProviderForm, false},
		{"/auth/github", handlers.LoginGithub, false},
		{"/auth/github/callback", handlers.LoginGithubCallback, false},
		{"/auth/google", handlers.LoginGoogle, false},
		{"/auth/google/callback", handlers.LoginGoogleCallback, false},
	}

	middleware := &middleware.Middleware{Handlers: handlers, Limiters: map[string]chan time.Time{}, BlockedSessions: map[string]time.Time{}}

	for _, route := range routes {
		mux.Handle(route.Path, middleware.DynamicMiddleware(route.Handler, route.RequireAuth))
	}

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads/"))))

	return middleware.StandardMiddleware(mux)
}
