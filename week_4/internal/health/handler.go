package health

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct {
	healthService *HealthService
}

func NewHealthHandler(healthService *HealthService) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
	}
}

func (h* HealthHandler) checkHealth(w http.ResponseWriter, r *http.Request) {
		response := h.healthService.healthCheck()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
}