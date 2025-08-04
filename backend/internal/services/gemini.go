package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/josephed37/FactCheck-AI/backend/internal/models"
	"github.com/josephed37/FactCheck-AI/backend/internal/search"
)

const geminiAPIURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key="

// GeminiService now holds a dependency on the TavilyService.
type GeminiService struct {
	SearchService *search.TavilyService
}

// FactCheck now orchestrates the RAG (Retrieval-Augmented Generation) process.
func (s *GeminiService) FactCheck(statement string) (*models.GeminiResponse, error) {
	// --- RAG Step 1: Retrieve ---
	searchResults, err := s.SearchService.Search(statement)
	if err != nil {
		return nil, fmt.Errorf("RAG search step failed: %w", err)
	}

	// --- RAG Step 2: Augment ---
	var contextBuilder strings.Builder
	contextBuilder.WriteString("Here is the real-time context from a web search:\n")
	for _, result := range searchResults {
		contextBuilder.WriteString(fmt.Sprintf("- Source: %s, Content: %s\n", result.URL, result.Content))
	}
	liveContext := contextBuilder.String()

	// Load API Key securely from environment
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	// Load the prompt template.
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filepath.Dir(filepath.Dir(b)))
	promptPath := filepath.Join(basepath, "prompts", "fact_check_prompt.txt")

	promptBytes, err := os.ReadFile(promptPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read prompt file at %s: %w", promptPath, err)
	}
	promptTemplate := string(promptBytes)
	augmentedPrompt := fmt.Sprintf(promptTemplate, liveContext, statement)

	// Construct the request body for the Gemini API
	requestBody, err := json.Marshal(map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{
						"text": augmentedPrompt,
					},
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create and send the HTTP POST request
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", geminiAPIURL+apiKey, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Gemini API: %w", err)
	}
	defer resp.Body.Close()

	// Read and parse the response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini API returned non-200 status: %s", string(responseBody))
	}

	var apiResponse struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(responseBody, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal outer API response: %w", err)
	}

	if len(apiResponse.Candidates) == 0 || len(apiResponse.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("invalid or empty response structure from Gemini")
	}

	geminiText := apiResponse.Candidates[0].Content.Parts[0].Text
	cleanedJSON := strings.TrimPrefix(geminiText, "```json")
	cleanedJSON = strings.TrimSuffix(cleanedJSON, "```")
	cleanedJSON = strings.TrimSpace(cleanedJSON)

	var geminiResp models.GeminiResponse
	if err := json.Unmarshal([]byte(cleanedJSON), &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal inner gemini response JSON (content: %s): %w", cleanedJSON, err)
	}

	// --- THIS IS THE CRITICAL ADDITION ---
	// We create a slice of our new Source model and populate it from the search results.
	var sources []models.Source
	for _, result := range searchResults {
		sources = append(sources, models.Source{
			Title: result.Title,
			URL:   result.URL,
		})
	}
	// We add the populated slice to our final response object before returning it.
	geminiResp.Sources = sources

	return &geminiResp, nil
}
