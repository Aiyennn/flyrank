package public

import (
	"github.com/go-chi/chi/v5"
)

type PublicRoute struct {
	handler *PublicHandler
}

func NewPublicRoute(handler *PublicHandler) *PublicRoute {
	return &PublicRoute{
		handler: handler,
	}
}

func (pr *PublicRoute) New() chi.Router {
	r := chi.NewRouter()
	r.Get("/info", pr.handler.GetInfo)
	return r
}
