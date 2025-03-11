package github

import (
	"errors"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/handler"
	"github.com/itelman/forum/internal/service/oauth"
	"github.com/itelman/forum/internal/service/oauth/domain"
	oauthApi "github.com/itelman/forum/pkg/oauth"
	"github.com/itelman/forum/pkg/sesm"
	"net/http"
)

type handlers struct {
	*handler.Handlers
	oauth     oauth.Service
	githubApi oauthApi.AuthApi
}

func NewHandlers(handler *handler.Handlers, oauth oauth.Service, api oauthApi.AuthApi) *handlers {
	return &handlers{handler, oauth, api}
}

func (h *handlers) RegisterMux(mux *http.ServeMux) {
	authRoutes := []dto.Route{
		{"/user/login/github", dto.GetMethod, h.login},
		{"/user/login/github/callback", dto.GetMethod, h.callback},
	}

	for _, route := range authRoutes {
		mux.Handle(route.Path, h.DynMiddleware.Chain(h.DynMiddleware.ForbidAuthenticatedUser(http.HandlerFunc(route.Handler)), route.Path, route.Methods))
	}
}

func (h *handlers) login(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, h.githubApi.GetAuthUri(), http.StatusFound)
}

func (h *handlers) callback(w http.ResponseWriter, r *http.Request) {
	req, err := oauth.DecodeLoginUserInput(r, h.githubApi)
	if err != nil {
		http.Redirect(w, r, "/user/login", http.StatusFound)
		return
	}

	resp, err := h.oauth.LoginUser(req.(*oauth.LoginUserInput))
	if errors.Is(err, domain.ErrOAuthUserNotFound) {
		http.Redirect(w, r, "/user/signup", http.StatusFound)
		return
	} else if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	//h.SesManager.DeleteActiveUserSession(resp.UserID)
	//http.SetCookie(w, dto.DeleteCookie(sesm.SessionId))

	sessionID, err := h.SesManager.CreateSession(resp.UserID)
	if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	http.SetCookie(w, dto.NewCookie(sesm.SessionId, sessionID))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
