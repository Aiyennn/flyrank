package auth

import (
	"github.com/go-chi/chi/v5"

	"week_4/internal/middleware"
)

type AuthRoute struct {
	authHandler *AuthHandler
}

func NewAuthRoute(authHandler *AuthHandler) *AuthRoute {
	return &AuthRoute{
		authHandler: authHandler,
	}
}

func (a *AuthRoute) New(authService *AuthService) chi.Router {
	r := chi.NewRouter()

	r.Post("/signup", a.authHandler.signUp)
	r.Post("/login", a.authHandler.login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(authService))
		r.Post("/logout", a.authHandler.logout)
	})

	return r
}

