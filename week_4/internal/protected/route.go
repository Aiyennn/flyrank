package protected

import (
	"github.com/go-chi/chi/v5"

	"week_4/internal/auth"
	"week_4/internal/middleware"
)

type ProtectedRoute struct {
	handler *ProtectedHandler
}

func NewProtectedRoute(handler *ProtectedHandler) *ProtectedRoute {
	return &ProtectedRoute{
		handler: handler,
	}
}

func (pr *ProtectedRoute) New(authService *auth.AuthService) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware(authService))
	r.Get("/profile", pr.handler.GetProfile)
	r.Get("/dashboard", pr.handler.GetDashboard)
	return r
}
