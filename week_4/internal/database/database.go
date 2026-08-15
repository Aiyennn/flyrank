package database

import (
	"fmt"

	supabase "github.com/supabase-community/supabase-go"
)

func NewClient(supabaseURL, supabaseKey string) (*supabase.Client, error) {
	if supabaseURL == "" || supabaseKey == "" {
		return nil, fmt.Errorf("SUPABASE_URL and SUPABASE_KEY must be set")
	}

	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		return nil, fmt.Errorf("create supabase client: %w", err)
	}

	return client, nil
}