package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"week_4/internal/health"
)

func New(db *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()

	healthService := health.NewHealthService()
	healthHandler := health.NewHealthHandler(healthService)
	healthRoute := health.NewHealthRoute(healthHandler)

	r.Mount("/health", healthRoute.New())

	return r
}