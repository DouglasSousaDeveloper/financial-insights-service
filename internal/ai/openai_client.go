package ai

import (
	"context"
	"fmt"

	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/domain"
	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/insight"
)

// OpenAIClient é a implementação concreta externa que bate nas APIs da OpenAI (Adapter)
type OpenAIClient struct {
	apiKey string
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		apiKey: apiKey,
	}
}

// GenerateFinancialInsight constrói um prompt otimizado baseado apenas nos agregados puros calculados no PG
func (c *OpenAIClient) GenerateFinancialInsight(ctx context.Context, summary *domain.FinancialSummary) (string, error) {
	// IMPORTANTE: Em um cenário real com "github.com/sashabaranov/go-openai",
	// criaríamos um cliente openai.NewClient(c.apiKey) aqui e enviaríamos a completion request.

	// Simulando a construção inteligente e minimalista do prompt, poupando tokens.
	prompt := fmt.Sprintf(
		"Você é um especialista financeiro. Elabore um curto insight baseado nisto: Rendas: %f | Despesas: %f",
		summary.TotalIncome, summary.TotalExpense,
	)
	_ = prompt // ignora aviso de variável não retornada (mock purposes)

	saldo := summary.TotalIncome - summary.TotalExpense
	respostaMockadaLLM := fmt.Sprintf("Olá! Com base nas suas finanças do período, você teve R$%.2f de rendas "+
		"e R$%.2f de despesas, resultando num saldo final de R$%.2f. Mantenha os bons hábitos!",
		summary.TotalIncome, summary.TotalExpense, saldo)

	return respostaMockadaLLM, nil
}

// Check para garantir que OpenAIClient atenda perfeitamente a interface abstracta insight.AIClient
var _ insight.AIClient = (*OpenAIClient)(nil)
