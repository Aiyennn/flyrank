package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	SupabaseURL string
	SupabaseKey string
}

func Load() (*Config, error) {
	if err := godotenv.Load("../.env"); err != nil {
		return nil, err
	}
	return &Config {
		DatabaseURL: os.Getenv("DATABASE_URL"),
		SupabaseURL: os.Getenv("SUPABASE_URL"),
		SupabaseKey: os.Getenv("SUPABASE_KEY"),
	}, nil
}