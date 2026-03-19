package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresDB inicializa o pool de conexões usando o pacote pgx.
// Nota: Em cenários reais, o "connectionString" viria do pacote de variáveis de ambiente.
func NewPostgresDB(ctx context.Context, connectionString string) (*pgxpool.Pool, error) {
	// ParseConfig faz o parse da string de conexão (ex: "postgres://user:pass@host:5432/db")
	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("não foi possível interpretar a configuração do banco: %w", err)
	}

	// Ajustes finos do pool
	poolConfig.MaxConns = 15 // Exemplo de controle de conexões simultâneas

	// Cria o Pool de instâncias em background
	db, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar no banco de dados: %w", err)
	}

	// Ping assegura que a conexão física está válida
	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping falhou: %w", err)
	}

	return db, nil
}
