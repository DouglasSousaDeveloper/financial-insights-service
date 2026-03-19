package transaction

import (
	"encoding/json"
	"net/http"

	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/domain"
)

// Handler estrutura a ponta de comunicação HTTP e recebe as dependências de Serviço.
// Em C# isso seria nosso "TransactionController".
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
	// Em Go nativo (sem libs), precisamos validar o verbo HTTP manualmente na mão
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var t domain.Transaction
	
	// A primitiva de Decoder do Go é super rápida para transformar JSON em Struct (Desserialização)
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Entrega a transação lida do JSON diretamente pro Serviço principal
	if err := h.service.ProcessTransaction(r.Context(), &t); err != nil {
		// Logaríamos o erro e devolveríamos 500 ou 400.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction created successfully"})
}
