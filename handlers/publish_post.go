package handlers

import (
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"dc-handler/services"
)

type PublishPostRequest struct {
	PostID       int    `json:"postId"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	FeaturedImageURL string `json:"featuredImageUrl"`
}

type PublishPostResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	PostID  int    `json:"postId"`
}

type PublishPostHandler struct {
	wpService *services.WordPressService
}

func NewPublishPostHandler(wpService *services.WordPressService) *PublishPostHandler {
	return &PublishPostHandler{wpService: wpService}
}

// PublishPost handles POST /api/publish-post
func (h *PublishPostHandler) PublishPost(c *gin.Context) {
	var req PublishPostRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	var featuredMediaID int

	// Upload featured image if provided
	if req.FeaturedImageURL != "" {
		filename := "featured-" + strconv.FormatInt(time.Now().UnixNano(), 10) + ".jpg"
		mediaID, err := h.wpService.UploadMedia(req.FeaturedImageURL, filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to upload featured image: " + err.Error(),
			})
			return
		}
		featuredMediaID = mediaID
	}

	// Update the post
	if err := h.wpService.UpdatePost(req.PostID, req.Title, req.Content, featuredMediaID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to publish post: " + err.Error(),
		})
		return
	}

	response := PublishPostResponse{
		Success: true,
		Message: "Post published successfully",
		PostID:  req.PostID,
	}

	c.JSON(http.StatusOK, response)
}
