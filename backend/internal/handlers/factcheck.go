package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/josephed37/FactCheck-AI/backend/internal/models"
	"github.com/josephed37/FactCheck-AI/backend/internal/services"
)

// FactCheckHandler is a struct that holds our services.
// This makes it easy to "inject" dependencies like our Gemini service.
type FactCheckHandler struct {
	GeminiService *services.GeminiService
}

// NewFactCheckHandler creates a new handler with its dependencies.
func NewFactCheckHandler(gs *services.GeminiService) *FactCheckHandler {
	return &FactCheckHandler{
		GeminiService: gs,
	}
}

// HandleFactCheck is the main function that processes the /fact-check request.
func (h *FactCheckHandler) HandleFactCheck(c *gin.Context) {
	// 1. Define a variable to hold the incoming request data.
	var req models.FactCheckRequest

	// 2. Bind the incoming JSON from the request body to our 'req' struct.
	// If the JSON is malformed or doesn't match our struct, this will fail.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return // Stop processing
	}

	// 3. Perform basic validation.
	if req.Statement == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Statement cannot be empty"})
		return
	}

	// 4. Call the FactCheck method from our Gemini service.
	// This is the core logic of our application.
	result, err := h.GeminiService.FactCheck(req.Statement)
	if err != nil {
		// If the service returns an error, send back a server error response.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 5. If everything is successful, send back the 200 OK response
	// with the result from the Gemini API.
	c.JSON(http.StatusOK, result)
}
