package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"dc-handler/config"
	"dc-handler/handlers"
	"dc-handler/services"
)

func main() {
	// Load configuration
	cfg := config.Load()

	fmt.Printf("🚀 Starting DC_Handler\n")
	fmt.Printf("   Environment: %s\n", cfg.Env)
	fmt.Printf("   Port: %s\n", cfg.Port)
	fmt.Printf("   WordPress: %s\n", cfg.WordPressURL)
	fmt.Printf("   Ollama: %s (%s)\n", cfg.OllamaURL, cfg.OllamaModel)
	fmt.Printf("\n")

	// Initialize services
	wpService := services.NewWordPressService(cfg)
	ollamaService := services.NewOllamaService(cfg)
	pixabayService := services.NewPixabayService(cfg)

	// Initialize handlers
	wpPostsHandler := handlers.NewWPPostsHandler(wpService)
	optimizeContentHandler := handlers.NewOptimizeContentHandler(ollamaService, pixabayService)
	publishPostHandler := handlers.NewPublishPostHandler(wpService)

	// Setup Gin router
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"env":     cfg.Env,
			"version": "1.0.0",
		})
	})

	// API Routes
	api := router.Group("/api")
	{
		// WordPress posts
		api.GET("/wp-posts", wpPostsHandler.GetPosts)

		// Optimize content
		api.POST("/optimize-content", optimizeContentHandler.OptimizeContent)

		// Publish post
		api.POST("/publish-post", publishPostHandler.PublishPost)
	}

	// Start server
	fmt.Printf("✓ Server running at http://localhost:%s\n", cfg.Port)
	fmt.Printf("✓ Health check: http://localhost:%s/health\n", cfg.Port)
	fmt.Printf("✓ API docs: see README.md\n\n")

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
