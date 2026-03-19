package transaction

import (
	"context"

	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/domain"
)

// Repository é a interface do lado do consumidor (Consumer-Side Interface).
// Quem injetar esta dependência (o adaptador do PostgreSQL) deverá implementar este contrato.
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

// ProcessTransaction aplica as regras de negócio e salva a transação.
func (s *Service) ProcessTransaction(ctx context.Context, t *domain.Transaction) error {
	// Aqui entrariam as validações de negócio reais. Exemplo:
	// if t.Amount <= 0 { return errors.New("amount must be positive") }
	
	// O serviço delega a persistência para a interface (Port)
	return s.repo.Save(ctx, t)
}
