package transaction

import (
	"context"

	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository implementa a interface transaction.Repository
type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Save(ctx context.Context, t *domain.Transaction) error {
	query := `
		INSERT INTO transactions (id, customer_id, amount, type, category, date)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	
	_, err := r.db.Exec(ctx, query, t.ID, t.CustomerID, t.Amount, t.Type, t.Category, t.Date)
	return err
}

var _ Repository = (*PostgresRepository)(nil)
