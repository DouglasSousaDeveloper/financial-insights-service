package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/database"
	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/insight"
	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/transaction"
)

func main() {
	// 1. Configura log estruturado nativo do Go (saída em JSON no console, ideal para CloudWatch/Datadog)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Contexto raiz de ciclo de vida da aplicação
	ctx := context.Background()

	logger.Info("Iniciando Financial Insights Service...")

	// 2. Tenta conectar ao Banco de Dados Postgres (usando pgxpool).
	// Em produção real, leríamos com os.Getenv("DATABASE_URL").
	dbString := "postgres://user:pass@localhost:5432/financialdb"
	dbPool, err := database.NewPostgresDB(ctx, dbString)
	if err != nil {
		logger.Error("Erro fatal ao conectar no banco", "error", err)
		// Em Go, ao invés de lançar uma Exception pro S.O., damos um log e encerramos o processo manualmente.
		os.Exit(1)
	}
	
	// A palavra-chave "defer" age como a interface "IDisposable" do C# (o antigo "using").
	// Garante que o método Close() do Pool seja executado assim que o main() for encerrado.
	defer dbPool.Close() 
	logger.Info("Conectado ao PostgreSQL com sucesso!")

	// 3. Injeção de Dependências Manual (Wiring Base)
	
	// -> Transaction Feature
	txRepo := transaction.NewPostgresRepository(dbPool)
	txService := transaction.NewService(txRepo)

	// -> Insight Feature
	insightRepo := insight.NewPostgresRepository(dbPool)
	// Por enquanto, comentamos o insightService porque precisamos da camada AIClient (que vamos construir)
	// mockAIClient := ... 
	// insightService := insight.NewService(insightRepo, mockAIClient)

	// (Passo Futuro) Aqui subiríamos o servidor HTTP (ex: echo, chi ou mux nativo) 
	// passando o txService e o insightService para os Handlers/Controllers.

	logger.Info("Injeção de dependências conectada:",
		"transaction_service", fmt.Sprintf("%T", txService),
		"insight_repo", fmt.Sprintf("%T", insightRepo),
	)

	// Impede que a main() encerre instantaneamente (é apenas um placeholder, 
	// já que o servidor HTTP que iremos adicionar no futuro possui um .ListenAndServe() bloqueante).
	select {}
}
