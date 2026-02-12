package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"dc-handler/services"
)

type WPPostsHandler struct {
	wpService *services.WordPressService
}

func NewWPPostsHandler(wpService *services.WordPressService) *WPPostsHandler {
	return &WPPostsHandler{wpService: wpService}
}

// GetPosts handles GET /api/wp-posts
func (h *WPPostsHandler) GetPosts(c *gin.Context) {
	status := c.DefaultQuery("status", "draft")
	perPageStr := c.DefaultQuery("per_page", "100")

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid per_page parameter",
		})
		return
	}

	posts, err := h.wpService.GetPosts(status, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"posts":   posts,
	})
}
