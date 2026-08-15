package protected

import (
	"github.com/go-chi/chi/v5"
)

type ProtectedRoute struct {
	handler *ProtectedHandler
}

func NewProtectedRoute(handler *ProtectedHandler) *ProtectedRoute {
	return &ProtectedRoute{
		handler: handler,
	}
}

func (pr *ProtectedRoute) New() chi.Router {
	r := chi.NewRouter()
	r.Get("/profile", pr.handler.GetProfile)
	return r
}
