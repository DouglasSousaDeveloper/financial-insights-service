package insight

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/domain"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GenerateInsight atende o POST /customers/{customer_id}/generate-insight
func (h *Handler) GenerateInsight(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	customerID := pathParts[2] 

	insight, err := h.service.GenerateInsight(r.Context(), customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(insight)
}

// GetInsights atende o GET /customers/{customer_id}/insights
func (h *Handler) GetInsights(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	customerID := pathParts[2]

	insights, err := h.service.GetInsights(r.Context(), customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if insights == nil {
		insights = []*domain.Insight{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(insights)
}
