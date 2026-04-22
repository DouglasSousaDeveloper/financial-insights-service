package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/ai"
	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/config"
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
	// Carrega Configurações de Ambiente
	cfg := config.LoadConfig()

	// Cria instância do Adapter de Banco de Dados (DB Connection)
	dbPool, err := database.NewPostgresDB(ctx, cfg.DatabaseURL)
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
	aiClient := ai.NewOpenAIClient(cfg.OpenAIApiKey)
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
	logger.Info("Iniciando servidor HTTP REST na porta " + cfg.AppPort + "...")
	if err := http.ListenAndServe(":"+cfg.AppPort, mux); err != nil {
		logger.Error("Falha crítica no servidor HTTP", "error", err)
		os.Exit(1)
	}
}
