package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/josephed37/FactCheck-AI/backend/config"
	"github.com/josephed37/FactCheck-AI/backend/internal/handlers"
	"github.com/josephed37/FactCheck-AI/backend/internal/services"
)

func main() {
	config.InitLogger()
	config.Log.Info("Starting the FactCheck-AI server...")

	router := gin.Default()

	// --- Middleware Configuration ---
	// We configure and apply the CORS middleware here. This must be done
	// before the routes are defined.
	router.Use(cors.New(cors.Config{
		// AllowOrigins specifies which origins are allowed to access the API.
		// For development, we allow our Streamlit app's address.
		AllowOrigins: []string{"http://localhost:8501"},

		// AllowMethods specifies which HTTP methods are allowed.
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},

		// AllowHeaders specifies which HTTP headers are allowed in requests.
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},

		// ExposeHeaders specifies which headers can be exposed as part of the response.
		ExposeHeaders: []string{"Content-Length"},

		// AllowCredentials allows cookies to be sent. Not needed now, but good practice.
		AllowCredentials: true,

		// MaxAge specifies how long the result of a preflight request can be cached.
		MaxAge: 12 * time.Hour,
	}))

	// --- Dependency Injection ---
	geminiService := &services.GeminiService{}
	factCheckHandler := handlers.NewFactCheckHandler(geminiService)

	// --- API Routes ---
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})
	router.POST("/v1/fact-check", factCheckHandler.HandleFactCheck)

	// --- Start Server ---
	config.Log.Info("Server is configured and running on port 8080")
	router.Run(":8080")
}
