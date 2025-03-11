package main

import (
	"errors"
	"fmt"
	"github.com/itelman/forum/internal/exception"
	"github.com/itelman/forum/internal/handler"
	activityHandlers "github.com/itelman/forum/internal/handler/activity"
	commentsHandlers "github.com/itelman/forum/internal/handler/comments"
	"github.com/itelman/forum/internal/handler/home"
	notificationsHandlers "github.com/itelman/forum/internal/handler/notifications"
	"github.com/itelman/forum/internal/handler/oauth/github"
	"github.com/itelman/forum/internal/handler/oauth/google"
	postsHandlers "github.com/itelman/forum/internal/handler/posts"
	commentReactionsHandlers "github.com/itelman/forum/internal/handler/reactions/comment_reactions"
	postReactionsHandlers "github.com/itelman/forum/internal/handler/reactions/post_reactions"
	usersHandlers "github.com/itelman/forum/internal/handler/users"
	"github.com/itelman/forum/internal/middleware/auth"
	"github.com/itelman/forum/internal/middleware/dynamic"
	"github.com/itelman/forum/internal/middleware/standard"
	"github.com/itelman/forum/internal/service/activity"
	"github.com/itelman/forum/internal/service/categories"
	"github.com/itelman/forum/internal/service/comment_reactions"
	"github.com/itelman/forum/internal/service/comments"
	"github.com/itelman/forum/internal/service/filters"
	"github.com/itelman/forum/internal/service/notifications"
	"github.com/itelman/forum/internal/service/oauth"
	"github.com/itelman/forum/internal/service/post_reactions"
	"github.com/itelman/forum/internal/service/posts"
	"github.com/itelman/forum/internal/service/users"
	"github.com/itelman/forum/pkg/templates"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer f.Close()

	conf := newConfig()
	deps, err := NewDependencies(
		WithSqlite(conf.Sqlite.DbDir, conf.Sqlite.MigrDir),
		WithGithubAuth(conf.Github.ClientSecret, conf.Github.ClientID, conf.ApiHost),
		WithGoogleAuth(conf.Google.ClientSecret, conf.Google.ClientID, conf.ApiHost),
		WithTemplateCache(conf.UI.TmplDir),
	)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer deps.Close()

	tmplRender := templates.NewTemplateRender(deps.templateCache, deps.sesManager)
	exceptionHandlers := exception.NewExceptions(errorLog, tmplRender)

	usersSvc := users.NewService(
		users.WithSqlite(deps.sqlite),
	)

	authMiddleware := auth.NewMiddleware(usersSvc, deps.sesManager, exceptionHandlers)
	dynamicMiddleware := dynamic.NewMiddleware(authMiddleware, deps.sesManager, exceptionHandlers)
	defaultHandlers := handler.NewHandlers(dynamicMiddleware, deps.sesManager, exceptionHandlers, tmplRender)

	postsSvc := posts.NewService(
		posts.WithSqlite(deps.sqlite),
	)

	commentsSvc := comments.NewService(
		comments.WithSqlite(deps.sqlite),
	)

	postReactionsSvc := post_reactions.NewService(
		post_reactions.WithSqlite(deps.sqlite),
	)

	commentReactionsSvc := comment_reactions.NewService(
		comment_reactions.WithSqlite(deps.sqlite),
	)

	categoriesSvc := categories.NewService(
		categories.WithSqlite(deps.sqlite),
	)

	filtersSvc := filters.NewService(
		filters.WithSqlite(deps.sqlite),
	)

	oauthSvc := oauth.NewService(
		oauth.WithSqlite(deps.sqlite),
	)

	notificationsSvc := notifications.NewService(
		notifications.WithSqlite(deps.sqlite),
	)

	activitySvc := activity.NewService(
		activity.WithSqlite(deps.sqlite),
	)

	mux := http.NewServeMux()

	home.NewHandlers(defaultHandlers, postsSvc, categoriesSvc, filtersSvc).RegisterMux(mux)
	usersHandlers.NewHandlers(defaultHandlers, usersSvc).RegisterMux(mux)
	postsHandlers.NewHandlers(defaultHandlers, postsSvc, commentsSvc, categoriesSvc, conf.PostImagesDir).RegisterMux(mux)
	commentsHandlers.NewHandlers(defaultHandlers, commentsSvc, postsSvc).RegisterMux(mux)
	postReactionsHandlers.NewHandlers(defaultHandlers, postReactionsSvc, postsSvc).RegisterMux(mux)
	commentReactionsHandlers.NewHandlers(defaultHandlers, commentReactionsSvc, commentsSvc).RegisterMux(mux)
	github.NewHandlers(defaultHandlers, oauthSvc, deps.githubAuth).RegisterMux(mux)
	google.NewHandlers(defaultHandlers, oauthSvc, deps.googleAuth).RegisterMux(mux)
	notificationsHandlers.NewHandlers(defaultHandlers, notificationsSvc).RegisterMux(mux)
	activityHandlers.NewHandlers(defaultHandlers, activitySvc).RegisterMux(mux)

	fileServer := http.FileServer(http.Dir(conf.UI.CSSDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(conf.PostImagesDir))))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", conf.Port),
		ErrorLog:     errorLog,
		Handler:      standard.NewMiddleware(exceptionHandlers, infoLog).Chain(mux),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		infoLog.Printf("Starting server on %s", conf.ApiHost)

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errorLog.Fatalf("Server error: %v", err)
		}
	}()

	sig := <-stop
	infoLog.Printf("Received shutdown signal: %v", sig)

	infoLog.Println("Server shutting down...")
}

/*
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
		},
	}
*/
