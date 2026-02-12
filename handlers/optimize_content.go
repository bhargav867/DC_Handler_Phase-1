package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dc-handler/services"
)

type OptimizeContentRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Excerpt string `json:"excerpt"`
}

type OptimizeContentResponse struct {
	OptimizedTitle   string `json:"optimizedTitle"`
	OptimizedContent string `json:"optimizedContent"`
	SuggestedImage   string `json:"suggestedImage"`
	ImageSource      string `json:"imageSource"`
}

type OptimizeContentHandler struct {
	ollamaService  *services.OllamaService
	pixabayService *services.PixabayService
}

func NewOptimizeContentHandler(ollamaService *services.OllamaService, pixabayService *services.PixabayService) *OptimizeContentHandler {
	return &OptimizeContentHandler{
		ollamaService:  ollamaService,
		pixabayService: pixabayService,
	}
}

// OptimizeContent handles POST /api/optimize-content
func (h *OptimizeContentHandler) OptimizeContent(c *gin.Context) {
	var req OptimizeContentRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Generate optimized content using Ollama
	optimizedContent, err := h.ollamaService.GenerateContent(req.Title, req.Content, req.Excerpt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to optimize content: " + err.Error(),
		})
		return
	}

	// Search for an image
	imageURL, err := h.pixabayService.SearchImage(req.Title)
	if err != nil {
		imageURL = ""
	}

	response := OptimizeContentResponse{
		OptimizedTitle:   req.Title + " - Enhanced",
		OptimizedContent: optimizedContent,
		SuggestedImage:   imageURL,
		ImageSource:      "Pixabay",
	}

	c.JSON(http.StatusOK, response)
}
