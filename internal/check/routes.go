package check

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", app.dynamicMiddleware(http.HandlerFunc(app.home)))
	mux.Handle("/results", app.dynamicMiddleware(http.HandlerFunc(app.results)))
	mux.Handle("/sneep/create", app.dynamicMiddleware(app.requireAuthenticatedUser(http.HandlerFunc(app.createSnippetForm))))
	mux.Handle("/like/", app.dynamicMiddleware(app.requireAuthenticatedUser(http.HandlerFunc(app.handleLike))))
	mux.Handle("/dislike/", app.dynamicMiddleware(app.requireAuthenticatedUser(http.HandlerFunc(app.handleDislike))))
	mux.Handle("/sneep/create/comment/", app.dynamicMiddleware(app.requireAuthenticatedUser(http.HandlerFunc(app.createComment))))
	mux.Handle("/sneep", app.dynamicMiddleware(http.HandlerFunc(app.snippet)))
	mux.Handle("/user/signup", app.dynamicMiddleware(http.HandlerFunc(app.signupUserForm)))
	mux.Handle("/user/login", app.dynamicMiddleware(http.HandlerFunc(app.loginUserForm)))
	mux.Handle("/user/logout", app.dynamicMiddleware(app.requireAuthenticatedUser(http.HandlerFunc(app.logoutUser))))
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	return app.standardMiddleware(mux)
}
