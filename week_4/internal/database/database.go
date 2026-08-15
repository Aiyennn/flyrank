package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, url string) (*pgxpool.Pool, error) {
    db, err := pgxpool.New(ctx, url)
    if err != nil {
        return nil, fmt.Errorf("create database pool: %w", err)
    }

    if err := db.Ping(ctx); err != nil {
        db.Close()
        return nil, fmt.Errorf("ping database: %w", err)
    }

    return db, nil
}