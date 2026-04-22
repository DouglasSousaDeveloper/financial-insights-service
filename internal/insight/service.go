package insight

import (
	"context"
	"time"

	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/domain"
	"github.com/google/uuid"
)

// Repository lida com a leitura e gravação no banco de dados.
type Repository interface {
	GetSummaryByCustomer(ctx context.Context, customerID string) (*domain.FinancialSummary, error)
	SaveInsight(ctx context.Context, insight *domain.Insight) error
	GetInsightsByCustomer(ctx context.Context, customerID string) ([]*domain.Insight, error)
}

// AIClient é a interface que conversa com o provedor de IA (OpenAI, Bedrock).
type AIClient interface {
	GenerateFinancialInsight(ctx context.Context, summary *domain.FinancialSummary) (string, error)
}

// Service orquestra o fluxo de geração de insights e contém as regras de domínio.
type Service struct {
	repo     Repository
	aiClient AIClient
}

// NewService cria uma nova instância do Service.
func NewService(repo Repository, aiClient AIClient) *Service {
	return &Service{
		repo:     repo,
		aiClient: aiClient,
	}
}

// GenerateInsight cria um novo insight para o usuário baseado no resumo financeiro.
func (s *Service) GenerateInsight(ctx context.Context, customerID string) (*domain.Insight, error) {
	summary, err := s.repo.GetSummaryByCustomer(ctx, customerID)
	if err != nil {
		return nil, err
	}

	content, err := s.aiClient.GenerateFinancialInsight(ctx, summary)
	if err != nil {
		return nil, err
	}

	insight := &domain.Insight{
		ID:         uuid.New().String(),
		CustomerID: customerID,
		Content:    content,
		CreatedAt:  time.Now(),
	}

	if err := s.repo.SaveInsight(ctx, insight); err != nil {
		return nil, err
	}

	return insight, nil
}

// GetInsights busca todos os insights vinculados a um determinado cliente.
func (s *Service) GetInsights(ctx context.Context, customerID string) ([]*domain.Insight, error) {
	return s.repo.GetInsightsByCustomer(ctx, customerID)
}
