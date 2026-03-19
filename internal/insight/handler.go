package insight

import (
	"encoding/json"
	"net/http"
	"strings"
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

	// Como optamos por NÃO instalar um router complexo (como Gorilla Mux ou Chi) ainda,
	// e para compatibilidade universal, extraímos o path variable "customer_id" cortando a string:
	// Path: /customers/123/generate-insight -> [ "", "customers", "123", "generate-insight" ]
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	customerID := pathParts[2] 

	// Chama o super Serviço que orquestra BD e IA (Observe como é assinado exatamente igual o Controller do C#)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Endpoint futuro: Aqui retornaremos os insights do array no DB para o cliente " + customerID,
	})
}
