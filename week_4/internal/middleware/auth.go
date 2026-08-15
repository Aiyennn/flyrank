package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/supabase-community/gotrue-go/types"
)

type contextKey string

const UserContextKey contextKey = "user"

type errorResponse struct {
	Error string `json:"error"`
}

type TokenVerifier interface {
	GetUserByToken(token string) (*types.UserResponse, error)
}

func AuthMiddleware(verifier TokenVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			const prefix = "Bearer "
			if authHeader == "" || len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(errorResponse{Error: "Access token required"})
				return
			}

			token := authHeader[len(prefix):]
			if token == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(errorResponse{Error: "Access token required"})
				return
			}

			user, err := verifier.GetUserByToken(token)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(errorResponse{Error: "Invalid or expired token"})
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
