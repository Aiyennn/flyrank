package protected

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/supabase-community/gotrue-go/types"
	"week_4/internal/auth"
	"week_4/internal/middleware"
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
	user, ok := r.Context().Value(middleware.UserContextKey).(*types.UserResponse)
	if !ok || user == nil {
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

func (h *ProtectedHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*types.UserResponse)
	if !ok || user == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errorResponse{Error: "Invalid or expired token"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Welcome to the dashboard, " + user.User.Email + "!",
	})
}
