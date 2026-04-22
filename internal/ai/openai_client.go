package ai

import (
	"context"
	"fmt"

	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/domain"
	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/insight"
	"github.com/sashabaranov/go-openai"
)

// OpenAIClient é a implementação concreta externa que bate nas APIs da OpenAI (Adapter)
type OpenAIClient struct {
	client *openai.Client
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		client: openai.NewClient(apiKey),
	}
}

// GenerateFinancialInsight chama a API da OpenAI para gerar insights baseado nos agregados
func (c *OpenAIClient) GenerateFinancialInsight(ctx context.Context, summary *domain.FinancialSummary) (string, error) {
	saldo := summary.TotalIncome - summary.TotalExpense
	prompt := fmt.Sprintf(
		"Você é um especialista financeiro rigoroso. Crie um parágrafo bem curto e amigável com um conselho para um cliente que teve:\n"+
		"Rendas Totais: R$%.2f\n"+
		"Despesas Totais: R$%.2f\n"+
		"Saldo final: R$%.2f\n",
		summary.TotalIncome, summary.TotalExpense, saldo,
	)

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		MaxTokens: 150,
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("falha ao comunicar com API da OpenAI: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("nenhuma resposta gerada pela OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}

// Check para garantir que OpenAIClient atenda perfeitamente a interface abstracta insight.AIClient
var _ insight.AIClient = (*OpenAIClient)(nil)
