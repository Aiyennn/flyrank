package auth

import (
	"github.com/go-chi/chi/v5"
)

type AuthRoute struct {
	authHandler *AuthHandler
}

func NewAuthRoute(authHandler *AuthHandler) *AuthRoute {
	return &AuthRoute{
		authHandler: authHandler,
	}
}

func (a *AuthRoute) New() chi.Router {
	r := chi.NewRouter()

	r.Post("/signup", a.authHandler.signUp)
	r.Post("/login", a.authHandler.login)

	return r
}

