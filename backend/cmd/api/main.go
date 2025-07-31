package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/josephed37/FactCheck-AI/backend/config"
	"github.com/josephed37/FactCheck-AI/backend/internal/database"
	"github.com/josephed37/FactCheck-AI/backend/internal/handlers"
	"github.com/josephed37/FactCheck-AI/backend/internal/services"
)

func main() {
	// --- Initialization ---
	// Set up our structured logger first.
	config.InitLogger()
	config.Log.Info("Logger initialized. Starting the FactCheck-AI server...")

	// NEW: Initialize the database connection.
	// We are telling our app to connect to a SQLite file named "factchecks.db"
	// inside the "data" directory. If this step fails, the application
	// cannot function correctly, so we log a fatal error and exit.
	if err := database.InitDB("./data/factchecks.db"); err != nil {
		config.Log.WithError(err).Fatal("Failed to initialize database connection")
	}
	config.Log.Info("Database connection initialized successfully")

	// Initialize the Gin router.
	router := gin.Default()

	// --- Middleware Configuration ---
	// Configure and apply the CORS middleware.
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8501"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// --- Dependency Injection ---
	// Create and connect all the major components of our application.
	geminiService := &services.GeminiService{}
	factCheckHandler := handlers.NewFactCheckHandler(geminiService)

	// --- API Routes ---
	// Define the specific URLs our server listens to.
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})
	router.POST("/v1/fact-check", factCheckHandler.HandleFactCheck)

	// Define the history endpoint.
	// - It listens for HTTP GET requests.
	// - It maps the URL to the `HandleGetHistory` method on our handler instance.
	router.GET("/v1/history", factCheckHandler.HandleGetHistory)

	// --- Start Server ---
	// Start listening for incoming requests on port 8080.
	config.Log.Info("Server is configured and running on port 8080")
	router.Run(":8080")
}
