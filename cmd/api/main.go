package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// 4. Configuração de Roteamento HTTP nativo
	mux := setupRoutes(txHandler, insightHandler)
	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: mux,
	}

	go func() {
		logger.Info("Iniciando servidor HTTP REST na porta " + cfg.AppPort + "...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Falha crítica no servidor HTTP", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Iniciando Graceful Shutdown...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		logger.Error("Erro ao encerrar o servidor", "error", err)
	}

	logger.Info("Servidor encerrado. Conexões finalizadas com segurança.")
}
