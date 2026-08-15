package health

import (
	"github.com/go-chi/chi/v5"
)

type HealthRoute struct {
	healthHandler *HealthHandler
}

func NewHealthRoute(healthHandler *HealthHandler) *HealthRoute {
	return &HealthRoute{
		healthHandler: healthHandler,
	}
}

func (h *HealthRoute) New() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.healthHandler.checkHealth)

	return r
}