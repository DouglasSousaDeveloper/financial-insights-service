package insight

import (
	"context"
	"fmt"

	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

// GetSummaryByCustomer busca os gastos e formata para o sumário financeiro.
func (r *PostgresRepository) GetSummaryByCustomer(ctx context.Context, customerID string) (*domain.FinancialSummary, error) {
	summaryQuery := `
		SELECT 
			COALESCE(SUM(amount) FILTER (WHERE type = 'INCOME'), 0) AS total_income,
			COALESCE(SUM(amount) FILTER (WHERE type = 'EXPENSE'), 0) AS total_expense
		FROM transactions
		WHERE customer_id = $1
	`
	
	summary := &domain.FinancialSummary{
		CustomerID:     customerID,
		CategorySpends: make(map[string]float64),
	}

	err := r.db.QueryRow(ctx, summaryQuery, customerID).Scan(&summary.TotalIncome, &summary.TotalExpense)
	if err != nil {
		return nil, fmt.Errorf("falha ao analisar os agrupamentos financeiros: %w", err)
	}

	categoryQuery := `
		SELECT category, SUM(amount)
		FROM transactions
		WHERE customer_id = $1 AND type = 'EXPENSE'
		GROUP BY category
	`
	rows, err := r.db.Query(ctx, categoryQuery, customerID)
	if err != nil {
		return nil, fmt.Errorf("falha ao realizar o agrupamento por categoria: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		var amount float64
		if err := rows.Scan(&category, &amount); err != nil {
			return nil, err
		}
		summary.CategorySpends[category] = amount
	}

	return summary, nil
}

func (r *PostgresRepository) SaveInsight(ctx context.Context, insight *domain.Insight) error {
	query := `
		INSERT INTO insights (id, customer_id, content, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(ctx, query, insight.ID, insight.CustomerID, insight.Content, insight.CreatedAt)
	return err
}

// GetInsightsByCustomer busca uma lista de insights de um usuário ordenados pela criação descendente.
func (r *PostgresRepository) GetInsightsByCustomer(ctx context.Context, customerID string) ([]*domain.Insight, error) {
	query := `
		SELECT id, customer_id, content, created_at
		FROM insights
		WHERE customer_id = $1
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar banco de dados de insights: %w", err)
	}
	defer rows.Close()

	var insights []*domain.Insight
	for rows.Next() {
		var insight domain.Insight
		if err := rows.Scan(&insight.ID, &insight.CustomerID, &insight.Content, &insight.CreatedAt); err != nil {
			return nil, err
		}
		insights = append(insights, &insight)
	}

	return insights, nil
}

var _ Repository = (*PostgresRepository)(nil)
