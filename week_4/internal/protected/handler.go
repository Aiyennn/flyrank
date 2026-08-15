package protected

import (
	"encoding/json"
	"net/http"
	"time"

	"week_4/internal/auth"
)

type ProtectedHandler struct {
	authService *auth.AuthService
}

func NewProtectedHandler(authService *auth.AuthService) *ProtectedHandler {
	return &ProtectedHandler{
		authService: authService,
	}
}

type errorResponse struct {
	Error string `json:"error"`
}

type profileResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *ProtectedHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.authService.GetUserByToken(token)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errorResponse{Error: "Invalid or expired token"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(profileResponse{
		ID:        user.User.ID.String(),
		Email:     user.User.Email,
		CreatedAt: user.User.CreatedAt,
	})
}
