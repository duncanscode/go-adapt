package main

import (
	"fmt"
	"go-adapt/internal/content"
	"go-adapt/internal/handler"
	"go-adapt/internal/llm"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found - using system environment variables")
	}

	// Setup
	apiKey := os.Getenv("ANTHROPIC_API_KEY")

	bank := content.NewStaticBank()
	llmClient := llm.NewLLMClient(apiKey)
	if apiKey != "" {
		llmClient = llm.NewLLMClient(apiKey)
		fmt.Println("LLM client initialized (LLM mode available)")
	} else {
		fmt.Println("ANTHROPIC_API_KEY not set - LLM mode disabled")
	}

	h := handler.NewHandler(bank, llmClient)

	// Configure Gin for production
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = "debug" // Default to debug for local development
	}
	gin.SetMode(mode)

	// Define routes
	r := gin.Default()

	// Configure trusted proxies (Apache on same server)
	// Use private IP for internal communication with Apache
	err = r.SetTrustedProxies([]string{"127.0.0.1", "10.124.0.2"})
	if err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// Enable CORS for API requests from frontend only
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		// Allow requests from your frontend domain or when proxied by Apache (no Origin header)
		if origin == "" || origin == "http://go-adapt.duncan.wiki" || origin == "https://go-adapt.duncan.wiki" {
			if origin != "" {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			}
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes only - frontend is served by Apache
	r.POST("/session/start", h.StartSession)
	r.GET("/session/question", h.GetNextQuestion)
	r.POST("/session/answer", h.SubmitAnswer)
	r.GET("/session/metrics", h.GetMetrics)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start API server (this blocks forever, handling requests)
	port := os.Getenv("PORT")
	if port == "" {
		port = "1234"
	}
	fmt.Printf("API server starting in %s mode on port %s (frontend served by Apache)\n", mode, port)
	r.Run(":" + port)
}
