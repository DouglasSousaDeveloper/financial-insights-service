package domain

import (
	"time"
)

// FinancialSummary representa o agregado de transações do cliente (criado no DB e mapeado aqui).
type FinancialSummary struct {
	CustomerID      string             `json:"customer_id"`
	TotalIncome     float64            `json:"total_income"`
	TotalExpense    float64            `json:"total_expense"`
	CategorySpends  map[string]float64 `json:"category_spends"`
	MonthlyVariance float64            `json:"monthly_variance"` // Ex: -0.05 para queda de 5% num mês
}

// Insight representa o texto final gerado pela IA.
type Insight struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}
