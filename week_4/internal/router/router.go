package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	supabase "github.com/supabase-community/supabase-go"

	"week_4/internal/auth"
	"week_4/internal/health"
)

func New(db *supabase.Client) http.Handler {
	r := chi.NewRouter()

	healthService := health.NewHealthService()
	healthHandler := health.NewHealthHandler(healthService)
	healthRoute := health.NewHealthRoute(healthHandler)

	r.Mount("/health", healthRoute.New())

	authService := auth.NewAuthService(db)
	authHandler := auth.NewAuthHandler(authService)
	authRoute := auth.NewAuthRoute(authHandler)

	r.Mount("/auth", authRoute.New())

	return r
}