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

// GetSummaryByCustomer faz pesadas operações de banco e formata para sumário financeiro do Domínio.
func (r *PostgresRepository) GetSummaryByCustomer(ctx context.Context, customerID string) (*domain.FinancialSummary, error) {
	// O banco que irá totalizar agrupando a entrada/saida inteira!
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

	// QueryRow é usado para retornar APENAS UMA linha (o consolidado dele)
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
	// Tratamos as categorias com Query, que retorna N linhas agrupadas.
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

var _ Repository = (*PostgresRepository)(nil)
