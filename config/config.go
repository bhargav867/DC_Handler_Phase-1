package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server config
	Port string
	Env  string

	// WordPress config
	WordPressURL      string
	WordPressUsername string
	WordPressPassword string

	// OpenAI config (will be replaced with Ollama)
	OpenAIAPIKey string
	OpenAIModel  string

	// Ollama config (for local LLM)
	OllamaURL   string
	OllamaModel string

	// Pixabay config
	PixabayAPIKey string
}

func Load() *Config {
	// Load .env file if it exists
	_ = godotenv.Load(".env")

	cfg := &Config{
		Port:              getEnv("PORT", "8080"),
		Env:               getEnv("ENV", "development"),
		WordPressURL:      getEnv("WORDPRESS_URL", "https://digitalchew.com"),
		WordPressUsername: getEnv("WORDPRESS_USERNAME", ""),
		WordPressPassword: getEnv("WORDPRESS_PASSWORD", ""),
		OpenAIAPIKey:      getEnv("OPENAI_API_KEY", ""),
		OpenAIModel:       getEnv("OPENAI_MODEL", "gpt-4o"),
		OllamaURL:         getEnv("OLLAMA_URL", "http://localhost:11434"),
		OllamaModel:       getEnv("OLLAMA_MODEL", "mistral"),
		PixabayAPIKey:     getEnv("PIXABAY_API_KEY", ""),
	}

	// Validate required config
	if cfg.WordPressURL == "" {
		fmt.Println("⚠️  Warning: WORDPRESS_URL not set")
	}

	return cfg
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
