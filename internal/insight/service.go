package insight

import (
	"context"
	"time"

	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/domain"
	"github.com/google/uuid"
)

// Repository lida com a leitura e gravação no banco de dados.
// O processamento pesado (agregação SUM, GROUP BY) acontece via SQL no banco, e é retornado no formato do summary.
type Repository interface {
	GetSummaryByCustomer(ctx context.Context, customerID string) (*domain.FinancialSummary, error)
	SaveInsight(ctx context.Context, insight *domain.Insight) error
}

// AIClient é a interface que conversa com o provedor de IA (OpenAI, Bedrock).
// Só enviamos o 'summary' (transações agregadas) para reduzir os custos (tokens).
type AIClient interface {
	GenerateFinancialInsight(ctx context.Context, summary *domain.FinancialSummary) (string, error)
}

// Service orquestra o fluxo de geração de insights e contém as regras de domínio.
type Service struct {
	repo     Repository
	aiClient AIClient
}

// NewService construtor responsável por injetar as dependências (Dependency Injection).
func NewService(repo Repository, aiClient AIClient) *Service {
	return &Service{
		repo:     repo,
		aiClient: aiClient,
	}
}

// GenerateInsight é o principal Caso de Uso (Use Case) do sistema.
func (s *Service) GenerateInsight(ctx context.Context, customerID string) (*domain.Insight, error) {
	// 1. Busca os dados agregados diretamente do banco de dados (o banco faz o trabalho duro "GROUP BY")
	summary, err := s.repo.GetSummaryByCustomer(ctx, customerID)
	if err != nil {
		return nil, err
	}

	// 2. Chama a IA generativa passando os dados agregados
	content, err := s.aiClient.GenerateFinancialInsight(ctx, summary)
	if err != nil {
		return nil, err
	}

	// 3. Monta o objeto de domínio final
	insight := &domain.Insight{
		ID:         uuid.New().String(),
		CustomerID: customerID,
		Content:    content,
		CreatedAt:  time.Now(),
	}

	// 4. Salva o insight gerado no DB
	if err := s.repo.SaveInsight(ctx, insight); err != nil {
		return nil, err
	}

	return insight, nil
}
