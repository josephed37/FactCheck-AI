// backend/internal/handlers/factcheck.go

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/josephed37/FactCheck-AI/backend/config"
	"github.com/josephed37/FactCheck-AI/backend/internal/database"
	"github.com/josephed37/FactCheck-AI/backend/internal/models"
	"github.com/josephed37/FactCheck-AI/backend/internal/services"
	"github.com/sirupsen/logrus"
)

// FactCheckHandler is a struct that holds the services this handler depends on.
type FactCheckHandler struct {
	GeminiService *services.GeminiService
}

// NewFactCheckHandler is a "constructor" function to create a new handler.
func NewFactCheckHandler(gs *services.GeminiService) *FactCheckHandler {
	return &FactCheckHandler{
		GeminiService: gs,
	}
}

// HandleFactCheck is the core method that processes the API request to check a statement.
func (h *FactCheckHandler) HandleFactCheck(c *gin.Context) {
	// ... (This function is unchanged)
	var req models.FactCheckRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		config.Log.WithError(err).Warn("Invalid JSON in request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body format"})
		return
	}

	if req.Statement == "" {
		config.Log.Warn("Request received with empty statement field")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Statement cannot be empty"})
		return
	}

	config.Log.WithFields(logrus.Fields{
		"client_ip":        c.ClientIP(),
		"statement_length": len(req.Statement),
	}).Info("Received new valid fact-check request")

	result, err := h.GeminiService.FactCheck(req.Statement)
	if err != nil {
		config.Log.WithError(err).Error("An error occurred in the GeminiService")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	if err := database.SaveFactCheck(req, *result); err != nil {
		config.Log.WithError(err).Warn("Failed to save fact-check record to the database")
	} else {
		config.Log.Info("Successfully saved fact-check record to the database")
	}

	c.JSON(http.StatusOK, result)
}

// HandleGetHistory is the method that processes requests to fetch all past fact-checks.
func (h *FactCheckHandler) HandleGetHistory(c *gin.Context) {
	config.Log.WithFields(logrus.Fields{
		"client_ip": c.ClientIP(),
	}).Info("Received request for fact-check history")

	history, err := database.GetFactCheckHistory()
	if err != nil {
		config.Log.WithError(err).Error("Failed to retrieve fact-check history from database")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve history"})
		return
	}

	// --- THE FIX ---
	// In Go, a 'nil' slice is marshalled to JSON 'null'. An empty slice is
	// marshalled to JSON '[]'. Our database function returns a nil slice if
	// there are no rows. We must check for this case and send a proper empty
	// JSON array to the client to prevent frontend errors.
	if history == nil {
		history = []models.FactCheckHistoryItem{}
	}

	c.JSON(http.StatusOK, history)
}
