package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/ai"
	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/database"
	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/insight"
	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/transaction"
)

func main() {
	// 1. Configuração de Log Base
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx := context.Background()
	logger.Info("Iniciando Financial Insights Service...")

	// 2. Cria instância do Adapter de Banco de Dados (DB Connection)
	dbString := "postgres://user:pass@localhost:5432/financialdb" 
	dbPool, err := database.NewPostgresDB(ctx, dbString)
	if err != nil {
		logger.Error("Erro fatal ao conectar no banco", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()
	logger.Info("Conectado ao PostgreSQL com sucesso!")

	// 3. Injeção de Dependências Manual - O "Composition Root" do Go!
	// ---------------------------------------------------------------
	
	// Transaction Core
	txRepo := transaction.NewPostgresRepository(dbPool)
	txService := transaction.NewService(txRepo)

	// Insight Core e AI Adapter
	insightRepo := insight.NewPostgresRepository(dbPool)
	aiClient := ai.NewOpenAIClient("fake-openai-token-apikey")
	insightService := insight.NewService(insightRepo, aiClient)

	// Injetamos os Services nos Controllers/Handlers
	txHandler := transaction.NewHandler(txService)
	insightHandler := insight.NewHandler(insightService)

	logger.Info("Injeção de dependências concluída",
		"transaction_service", fmt.Sprintf("%T", txService),
		"insight_service", fmt.Sprintf("%T", insightService),
	)

	// 4. Configuração de Roteamento HTTP nativo ("Servidor Web" da biblioteca raiz do Go)
	mux := http.NewServeMux()
	
	mux.HandleFunc("/transactions", txHandler.CreateTransaction)
	
	mux.HandleFunc("/customers/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/generate-insight") && r.Method == http.MethodPost {
			insightHandler.GenerateInsight(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/insights") && r.Method == http.MethodGet {
			insightHandler.GetInsights(w, r)
			return
		}
		http.NotFound(w, r)
	})

	// 5. Inicia o Servidor (Essa chamada é bloqueante, igual 'app.Run()' em C# minimal APIs)
	logger.Info("Iniciando servidor HTTP REST na porta :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		logger.Error("Falha crítica no servidor HTTP", "error", err)
		os.Exit(1)
	}
}
