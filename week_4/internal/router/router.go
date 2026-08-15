package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	supabase "github.com/supabase-community/supabase-go"

	"week_4/internal/auth"
	"week_4/internal/docs"
	"week_4/internal/health"
	"week_4/internal/protected"
	"week_4/internal/public"
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

	r.Mount("/auth", authRoute.New(authService))

	publicHandler := public.NewPublicHandler()
	publicRoute := public.NewPublicRoute(publicHandler)
	r.Mount("/public", publicRoute.New())

	protectedHandler := protected.NewProtectedHandler(authService)
	protectedRoute := protected.NewProtectedRoute(protectedHandler)
	r.Mount("/protected", protectedRoute.New(authService))

	docsRoute := docs.NewDocsRoute()
	r.Mount("/docs", docsRoute.New())
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	})

	return r
}