package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	SupabaseURL string
	SupabaseKey string
	Port        string
}

func Load() (*Config, error) {
	_ = godotenv.Load("../.env")

	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		SupabaseURL: os.Getenv("SUPABASE_URL"),
		SupabaseKey: os.Getenv("SUPABASE_KEY"),
		Port:        os.Getenv("PORT"),
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return cfg, nil
}
