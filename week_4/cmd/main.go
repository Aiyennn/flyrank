package main

import (
	"log"
	"net/http"

	"week_4/internal/config"
	"week_4/internal/database"
	"week_4/internal/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	db, err := database.NewClient(cfg.SupabaseURL, cfg.SupabaseKey)
	if err != nil {
		log.Fatal("failed to create supabase client:", err)
	}

	log.Println("Supabase client initialized")

	r := router.New(db)

	addr := ":" + cfg.Port
	log.Printf("starting server on %s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
