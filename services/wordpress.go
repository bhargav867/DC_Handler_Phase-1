package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"dc-handler/config"
	"dc-handler/utils"
)

type WordPressPost struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	Excerpt      string `json:"excerpt"`
	Status       string `json:"status"`
	FeaturedMedia int   `json:"featured_media"`
	Link         string `json:"link"`
}

type WordPressService struct {
	cfg *config.Config
}

func NewWordPressService(cfg *config.Config) *WordPressService {
	return &WordPressService{cfg: cfg}
}

// GetPosts fetches posts from WordPress
func (ws *WordPressService) GetPosts(status string, perPage int) ([]WordPressPost, error) {
	url := fmt.Sprintf("%s/wp-json/wp/v2/posts?status=%s&per_page=%d&_embed",
		ws.cfg.WordPressURL, status, perPage)

	utils.LogRequest("GET", url, ws.cfg.WordPressUsername)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication
	if ws.cfg.WordPressUsername != "" && ws.cfg.WordPressPassword != "" {
		req.Header.Add("Authorization", utils.CreateBasicAuth(
			ws.cfg.WordPressUsername,
			ws.cfg.WordPressPassword,
		))
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.LogError("WP.GetPosts", err.Error())
		return nil, fmt.Errorf("failed to fetch posts: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		utils.LogError("WP.GetPosts", fmt.Sprintf("Status %d", resp.StatusCode))
		return nil, fmt.Errorf("WordPress API error (%d): %s", resp.StatusCode, string(body))
	}

	var posts []WordPressPost
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		utils.LogError("WP.GetPosts", err.Error())
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	utils.LogSuccess("WP.GetPosts", fmt.Sprintf("Retrieved %d posts", len(posts)))
	return posts, nil
}

// UpdatePost updates a WordPress post
func (ws *WordPressService) UpdatePost(postID int, title, content string, featuredMediaID int) error {
	url := fmt.Sprintf("%s/wp-json/wp/v2/posts/%d", ws.cfg.WordPressURL, postID)

	utils.LogRequest("POST", url, ws.cfg.WordPressUsername)

	payload := map[string]interface{}{
		"title":          title,
		"content":        content,
		"featured_media": featuredMediaID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", utils.CreateBasicAuth(
		ws.cfg.WordPressUsername,
		ws.cfg.WordPressPassword,
	))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.LogError("WP.UpdatePost", err.Error())
		return fmt.Errorf("failed to update post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		utils.LogError("WP.UpdatePost", fmt.Sprintf("Status %d", resp.StatusCode))
		return fmt.Errorf("WordPress API error (%d): %s", resp.StatusCode, string(body))
	}

	utils.LogSuccess("WP.UpdatePost", fmt.Sprintf("Post %d updated", postID))
	return nil
}

// UploadMedia uploads an image to WordPress media library
func (ws *WordPressService) UploadMedia(imageURL, filename string) (int, error) {
	utils.LogRequest("POST", "media-upload", filename)

	// Fetch the image
	resp, err := http.Get(imageURL)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to fetch image: status %d", resp.StatusCode)
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read image data: %w", err)
	}

	// Upload to WordPress
	url := fmt.Sprintf("%s/wp-json/wp/v2/media", ws.cfg.WordPressURL)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(imageData))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", utils.CreateBasicAuth(
		ws.cfg.WordPressUsername,
		ws.cfg.WordPressPassword,
	))
	req.Header.Add("Content-Type", "image/jpeg")
	req.Header.Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	client := &http.Client{}
	uploadResp, err := client.Do(req)
	if err != nil {
		utils.LogError("WP.UploadMedia", err.Error())
		return 0, fmt.Errorf("failed to upload media: %w", err)
	}
	defer uploadResp.Body.Close()

	if uploadResp.StatusCode != http.StatusCreated && uploadResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(uploadResp.Body)
		utils.LogError("WP.UploadMedia", fmt.Sprintf("Status %d", uploadResp.StatusCode))
		return 0, fmt.Errorf("WordPress media upload error (%d): %s", uploadResp.StatusCode, string(body))
	}

	var mediaResp map[string]interface{}
	if err := json.NewDecoder(uploadResp.Body).Decode(&mediaResp); err != nil {
		return 0, fmt.Errorf("failed to decode media response: %w", err)
	}

	mediaID := int(mediaResp["id"].(float64))
	utils.LogSuccess("WP.UploadMedia", fmt.Sprintf("Media ID: %d", mediaID))
	return mediaID, nil
}
