package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config representa as variáveis de ambiente necessárias para o projeto rodar.
type Config struct {
	DatabaseURL  string
	OpenAIApiKey string
	AppPort      string
}

// LoadConfig tenta carregar as varíaves do arquivo ".env" caso esteja rodando localmente,
// e preenche os valores default caso estejam ausentes do sistema operacional.
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: arquivo .env não encontrado. Confiando em variáveis pré-estabelecidas no SO")
	}

	return &Config{
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/financialdb"),
		OpenAIApiKey: getEnv("OPENAI_API_KEY", "fake-openai-token-apikey"),
		AppPort:      getEnv("APP_PORT", "8080"),
	}
}

// getEnv lê o valor de uma var de ambiente; se não existir, retorna um fallback
func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}
