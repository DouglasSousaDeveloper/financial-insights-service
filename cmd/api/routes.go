package main

import (
	"net/http"
	"strings"

	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/insight"
	"github.com/DouglasSousaDeveloper/financial-insights-service/internal/transaction"
)

func setupRoutes(txHandler *transaction.Handler, insightHandler *insight.Handler) *http.ServeMux {
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

	return mux
}
