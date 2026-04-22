package transaction

import (
	"context"
	"errors"
	"strings"

	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/domain"
	"github.com/google/uuid"
)

// Repository é a interface do lado do consumidor (Consumer-Side Interface).
type Repository interface {
	Save(ctx context.Context, t *domain.Transaction) error
}

// Service contém as regras de negócio para transações.
type Service struct {
	repo Repository
}

// NewService cria uma nova instância do serviço (Design Pattern: Factory/Constructor).
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) ProcessTransaction(ctx context.Context, t *domain.Transaction) error {
	if strings.TrimSpace(t.CustomerID) == "" {
		return errors.New("o customer_id é obrigatório")
	}
	if t.Amount <= 0 {
		return errors.New("o valor da transação deve ser maior que zero")
	}
	if t.Type != domain.TransactionTypeIncome && t.Type != domain.TransactionTypeExpense {
		return errors.New("o tipo da transação deve ser INCOME ou EXPENSE")
	}
	if strings.TrimSpace(t.Category) == "" {
		return errors.New("a categoria é obrigatória")
	}

	t.ID = uuid.New().String()

	return s.repo.Save(ctx, t)
}
