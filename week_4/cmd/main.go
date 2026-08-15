package main

import (
	"context"
	"log"

	"week_4/internal/config"
	"week_4/internal/database"
	"week_4/internal/health"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	db, err := database.Connect(ctx, cfg.SupabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("Connected to Supabase")

	healthService := health.NewHealthService()
	healthHandler := health.NewHealthHandler(healthService)
	healthRoute := health.NewHealthRoute(healthHandler)

	r := healthRoute.New()

	log.Println("starting server on :8080")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
