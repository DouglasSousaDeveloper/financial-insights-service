package domain

import (
	"time"
)

// TransactionType define se a transação é entrada ou saída.
type TransactionType string

const (
	TransactionTypeIncome  TransactionType = "INCOME"
	TransactionTypeExpense TransactionType = "EXPENSE"
)

// Transaction representa uma transação crua efetuada pelo cliente.
type Transaction struct {
	ID         string          `json:"id"`
	CustomerID string          `json:"customer_id"`
	Amount     float64         `json:"amount"` // Em produção real, prefira usar int64 (centavos) ou uma lib como shopspring/decimal
	Type       TransactionType `json:"type"`
	Category   string          `json:"category"`
	Date       time.Time       `json:"date"`
}
