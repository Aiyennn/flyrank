package public

import (
	"encoding/json"
	"net/http"
)

type PublicHandler struct {
}

func NewPublicHandler() *PublicHandler {
	return &PublicHandler{}
}

type infoResponse struct {
	Message string `json:"message"`
}

func (h *PublicHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(infoResponse{
		Message: "Welcome stranger! This info is public.",
	})
}
