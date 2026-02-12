package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"dc-handler/config"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

type OllamaService struct {
	cfg *config.Config
}

func NewOllamaService(cfg *config.Config) *OllamaService {
	return &OllamaService{cfg: cfg}
}

// GenerateContent generates optimized content using Ollama
func (os *OllamaService) GenerateContent(title, content, excerpt string) (string, error) {
	prompt := fmt.Sprintf(`You are an expert content optimizer. Please optimize the following blog post content:

Title: %s
Content: %s
Excerpt: %s

Provide an optimized version with better SEO, readability, and engagement. Return only the optimized content without explanations.`, title, content, excerpt)

	reqBody := OllamaRequest{
		Model:  os.cfg.OllamaModel,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/generate", os.cfg.OllamaURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API error (%d): %s", resp.StatusCode, string(body))
	}

	var olResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&olResp); err != nil {
		return "", fmt.Errorf("failed to decode Ollama response: %w", err)
	}

	return olResp.Response, nil
}
