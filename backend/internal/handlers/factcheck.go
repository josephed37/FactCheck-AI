package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/josephed37/FactCheck-AI/backend/config"
	"github.com/josephed37/FactCheck-AI/backend/internal/models"
	"github.com/josephed37/FactCheck-AI/backend/internal/services"
	"github.com/sirupsen/logrus"
)

// FactCheckHandler is a struct that holds the services this handler depends on.
// Analogy: Think of this as the "manager" of a department. The manager doesn't
// do the core work itself but knows which "employee" (service) to delegate tasks to.
// This design makes our code modular and much easier to test.
type FactCheckHandler struct {
	GeminiService *services.GeminiService
}

// NewFactCheckHandler is a "constructor" function. It's a clean and standard way
// to create a new FactCheckHandler while ensuring it receives all the dependencies
// (like the GeminiService) it needs to function correctly.
func NewFactCheckHandler(gs *services.GeminiService) *FactCheckHandler {
	return &FactCheckHandler{
		GeminiService: gs,
	}
}

// HandleFactCheck is the core method that processes the API request. It's the
// function that gets executed when a request hits `POST /v1/fact-check`.
func (h *FactCheckHandler) HandleFactCheck(c *gin.Context) {
	// 1. Prepare a variable to hold the incoming data.
	// We create an empty `FactCheckRequest` struct that we will populate.
	var req models.FactCheckRequest

	// 2. Bind and Validate the Request Body.
	// `c.ShouldBindJSON` is a powerful Gin function that does two things:
	//    a. It reads the JSON data from the incoming request.
	//    b. It validates that the JSON has the correct fields and types to fit
	//       into our `req` struct. This prevents malformed requests.
	if err := c.ShouldBindJSON(&req); err != nil {
		// If binding fails, it means the client sent bad data. We log a warning
		// and send a `400 Bad Request` response, which is the correct HTTP status code.
		config.Log.WithError(err).Warn("Invalid JSON in request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body format"})
		return // We must return here to stop further execution.
	}

	// 3. Perform Business Logic Validation.
	// Even if the JSON is well-formed, we need to check if the data makes sense.
	// Here, we ensure the `statement` field is not just an empty string.
	if req.Statement == "" {
		config.Log.Warn("Request received with empty statement field")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Statement cannot be empty"})
		return
	}

	// 4. Log the valid incoming request for monitoring and debugging purposes.
	// We use structured logging (`WithFields`) to attach key-value pairs to the
	// log entry, making it easy to search and filter logs later.
	config.Log.WithFields(logrus.Fields{
		"client_ip":        c.ClientIP(),
		"statement_length": len(req.Statement),
	}).Info("Received new valid fact-check request")

	// 5. Delegate to the Service Layer.
	// The handler's job is done with the HTTP logic. It now calls the `FactCheck`
	// method on the `GeminiService` to perform the actual work of talking to the AI.
	// This separation of concerns is key to a clean architecture.
	result, err := h.GeminiService.FactCheck(req.Statement)
	if err != nil {
		// If the service layer returns an error (e.g., the Gemini API is down),
		// we log it as a critical error and return a `500 Internal Server Error`.
		// This tells the client that something went wrong on our end.
		config.Log.WithError(err).Error("An error occurred in the GeminiService")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	// 6. Send the Success Response.
	// If no errors occurred, we send a `200 OK` status code along with the
	// `result` we received from the service.
	c.JSON(http.StatusOK, result)
}
