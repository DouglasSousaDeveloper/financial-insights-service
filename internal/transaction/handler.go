package transaction

import (
	"encoding/json"
	"net/http"

	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/domain"
)

// Handler estrutura a ponta de comunicação HTTP e recebe as dependências de Serviço.
type Handler struct {
	service *Service
}

// NewHandler cria uma nova instância do Controller e injeta o service de negócio.
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateTransaction mapeia a rota POST /transactions
func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var t domain.Transaction
	
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.ProcessTransaction(r.Context(), &t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction created successfully"})
}
