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

	// Define routes
	r := gin.Default()

	// API routes
	r.POST("/session/start", h.StartSession)
	r.GET("/session/question", h.GetNextQuestion)
	r.POST("/session/answer", h.SubmitAnswer)

	// Serve static frontend files (must come after API routes)
	r.Static("/static", "./frontend")
	r.StaticFile("/", "./frontend/index.html")

	// Start server (this blocks forever, handling requests)
	fmt.Println("server starting on http://localhost:8080")
	r.Run(":1234")
}