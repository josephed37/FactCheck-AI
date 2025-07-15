// backend/cmd/api/main.go

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/josephed37/FactCheck-AI/backend/internal/handlers"
	"github.com/josephed37/FactCheck-AI/backend/internal/services"
)

func main() {
	router := gin.Default()

	// --- Dependency Injection Setup ---
	// 1. Create an instance of our GeminiService.
	geminiService := &services.GeminiService{}

	// 2. Create an instance of our FactCheckHandler, passing in the service.
	factCheckHandler := handlers.NewFactCheckHandler(geminiService)

	// --- API Routes ---
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// 3. Define the main API endpoint.
	// It will respond to POST requests at /v1/fact-check
	// and use the HandleFactCheck method from our handler.
	router.POST("/v1/fact-check", factCheckHandler.HandleFactCheck)

	// --- Start Server ---
	router.Run(":8080")
}
