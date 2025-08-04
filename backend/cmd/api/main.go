package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/josephed37/FactCheck-AI/backend/config"
	"github.com/josephed37/FactCheck-AI/backend/internal/database"
	"github.com/josephed37/FactCheck-AI/backend/internal/handlers"
	"github.com/josephed37/FactCheck-AI/backend/internal/search" // Import the new search package
	"github.com/josephed37/FactCheck-AI/backend/internal/services"
)

func main() {
	config.InitLogger()
	config.Log.Info("Logger initialized. Starting the FactCheck-AI server...")

	if err := database.InitDB("./data/factchecks.db"); err != nil {
		config.Log.WithError(err).Fatal("Failed to initialize database connection")
	}
	config.Log.Info("Database connection initialized successfully")

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8501"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// --- UPDATED: Dependency Injection for RAG ---
	// 1. Create an instance of our new TavilyService.
	tavilyService := &search.TavilyService{}

	// 2. Create the GeminiService and pass the TavilyService to it.
	geminiService := &services.GeminiService{
		SearchService: tavilyService,
	}

	// 3. Create the handler, passing in the fully-equipped GeminiService.
	factCheckHandler := handlers.NewFactCheckHandler(geminiService)

	// --- API Routes ---
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})
	router.POST("/v1/fact-check", factCheckHandler.HandleFactCheck)
	router.GET("/v1/history", factCheckHandler.HandleGetHistory)

	// --- Start Server ---
	config.Log.Info("Server is configured and running on port 8080")
	router.Run(":8080")
}
