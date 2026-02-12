package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"dc-handler/config"
)

type PixabayImage struct {
	URL              string `json:"largeImageURL"`
	WebformatURL     string `json:"webformatURL"`
}

type PixabayResponse struct {
	Hits []PixabayImage `json:"hits"`
}

type PixabayService struct {
	cfg *config.Config
}

func NewPixabayService(cfg *config.Config) *PixabayService {
	return &PixabayService{cfg: cfg}
}

// SearchImage searches for an image on Pixabay
func (ps *PixabayService) SearchImage(query string) (string, error) {
	if ps.cfg.PixabayAPIKey == "" {
		return "", fmt.Errorf("pixabay API key not configured")
	}

	// Clean up query
	query = strings.TrimSpace(query)
	query = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(query, "")

	if len(query) < 3 {
		return "", fmt.Errorf("query too short")
	}

	if len(query) > 100 {
		query = query[:100]
	}

	// Build URL
	params := url.Values{}
	params.Add("key", ps.cfg.PixabayAPIKey)
	params.Add("q", query)
	params.Add("image_type", "photo")
	params.Add("orientation", "horizontal")
	params.Add("per_page", "3")
	params.Add("safesearch", "true")

	searchURL := fmt.Sprintf("https://pixabay.com/api/?%s", params.Encode())

	resp, err := http.Get(searchURL)
	if err != nil {
		return "", fmt.Errorf("failed to search Pixabay: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Pixabay API error (%d): %s", resp.StatusCode, string(body))
	}

	var pixResp PixabayResponse
	if err := json.NewDecoder(resp.Body).Decode(&pixResp); err != nil {
		return "", fmt.Errorf("failed to decode Pixabay response: %w", err)
	}

	if len(pixResp.Hits) == 0 {
		return "", fmt.Errorf("no images found")
	}

	imageURL := pixResp.Hits[0].URL
	if imageURL == "" {
		imageURL = pixResp.Hits[0].WebformatURL
	}

	return imageURL, nil
}
