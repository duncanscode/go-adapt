package main

import (
	"embed"
	"fmt"
	"go-adapt/internal/content"
	"go-adapt/internal/handler"
	"go-adapt/internal/llm"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

//go:embed frontend
var frontendFS embed.FS

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

	// Configure trusted proxies (nginx on same server)
	// Use private IP for internal communication with nginx
	err = r.SetTrustedProxies([]string{"127.0.0.1", "10.124.0.2"})
	if err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// API routes
	r.POST("/session/start", h.StartSession)
	r.GET("/session/question", h.GetNextQuestion)
	r.POST("/session/answer", h.SubmitAnswer)
	r.GET("/session/metrics", h.GetMetrics)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Serve embedded frontend files (must come after API routes)
	frontendSubFS, err := fs.Sub(frontendFS, "frontend")
	if err != nil {
		log.Fatalf("Failed to create frontend sub-filesystem: %v", err)
	}
	r.StaticFS("/static", http.FS(frontendSubFS))
	r.GET("/", func(c *gin.Context) {
		c.FileFromFS("index.html", http.FS(frontendSubFS))
	})

	// Start server (this blocks forever, handling requests)
	port := os.Getenv("PORT")
	if port == "" {
		port = "1234"
	}
	fmt.Printf("Server starting in %s mode on port %s\n", mode, port)
	r.Run(":" + port)
}
