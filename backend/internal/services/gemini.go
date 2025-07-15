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
	"time"

	"github.com/josephed37/FactCheck-AIinternal/models"
)

const geminiAPIURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key="

// GeminiService encapsulates the logic for interacting with the Google Gemini API.
type GeminiService struct{}

// FactCheck sends a statement to the Gemini API for analysis.
func (s *GeminiService) FactCheck(statement string) (*models.GeminiResponse, error) {
	// 1. Load API Key securely from environment
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	// 2. Load the prompt template from the file system
	// Note: This path is relative to the project root where the binary will be run.
	promptPath, err := filepath.Abs("../../prompts/fact_check_prompt.txt")
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path: %w", err)
	}

	promptBytes, err := os.ReadFile(promptPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read prompt template file: %w", err)
	}
	promptTemplate := string(promptBytes)
	finalPrompt := fmt.Sprintf(promptTemplate, statement)

	// 3. Construct the request body for the Gemini API
	// This structure is specific to the Gemini API's requirements.
	requestBody, err := json.Marshal(map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{
						"text": finalPrompt,
					},
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// 4. Create and send the HTTP POST request
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

	// 5. Read and parse the response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini API returned non-200 status: %s", string(responseBody))
	}

	// The actual content is nested inside the API response.
	// We need to unmarshal the outer structure first.
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

	// The actual JSON we want is a string inside the 'Text' field.
	geminiText := apiResponse.Candidates[0].Content.Parts[0].Text

	var geminiResp models.GeminiResponse
	if err := json.Unmarshal([]byte(geminiText), &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal inner gemini response JSON: %w", err)
	}

	return &geminiResp, nil
}