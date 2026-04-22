package transaction_test

import (
	"context"
	"strings"
	"testing"

	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/domain"
	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/transaction"
)

type mockTransactionRepo struct {
	SaveFunc func(ctx context.Context, t *domain.Transaction) error
}

func (m *mockTransactionRepo) Save(ctx context.Context, t *domain.Transaction) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, t)
	}
	return nil
}

func TestProcessTransaction(t *testing.T) {
	// Table-driven tests (Padrão ouro em Go)
	tests := []struct {
		name          string
		inputTx       *domain.Transaction
		expectedError string
	}{
		{
			name: "Sucesso - Transacao Valida",
			inputTx: &domain.Transaction{
				CustomerID: "CUST-1",
				Amount:     150.0,
				Type:       domain.TransactionTypeIncome,
				Category:   "Salário",
			},
			expectedError: "",
		},
		{
			name: "Erro - CustomerID vazio",
			inputTx: &domain.Transaction{
				CustomerID: "",
				Amount:     150.0,
				Type:       domain.TransactionTypeIncome,
				Category:   "Salário",
			},
			expectedError: "customer_id é obrigatório",
		},
		{
			name: "Erro - Valor negativo ou zero",
			inputTx: &domain.Transaction{
				CustomerID: "CUST-1",
				Amount:     -50.0,
				Type:       domain.TransactionTypeIncome,
				Category:   "Salário",
			},
			expectedError: "valor da transação deve ser maior que zero",
		},
		{
			name: "Erro - Tipo invalido",
			inputTx: &domain.Transaction{
				CustomerID: "CUST-1",
				Amount:     150.0,
				Type:       "INVALID_TYPE",
				Category:   "Salário",
			},
			expectedError: "tipo da transação deve ser INCOME ou EXPENSE",
		},
		{
			name: "Erro - Categoria vazia",
			inputTx: &domain.Transaction{
				CustomerID: "CUST-1",
				Amount:     150.0,
				Type:       domain.TransactionTypeExpense,
				Category:   "   ", // Espaços em branco caem no TrimSpace
			},
			expectedError: "categoria é obrigatória",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			repo := &mockTransactionRepo{}
			service := transaction.NewService(repo)

			// Act
			err := service.ProcessTransaction(context.Background(), tt.inputTx)

			// Assert
			if tt.expectedError == "" {
				if err != nil {
					t.Fatalf("esperava sucesso, mas deu erro: %v", err)
				}
				if tt.inputTx.ID == "" {
					t.Error("esperava que um UUID fosse gerado para a transação, encontrou string vazia")
				}
			} else {
				if err == nil {
					t.Fatalf("esperava erro contendo '%s', mas obteve sucesso", tt.expectedError)
				}
				if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("esperava erro '%s', obteve '%v'", tt.expectedError, err)
				}
			}
		})
	}
}
