// backend/internal/services/gemini.go

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
)

// Corrected: This is now a simple string constant.
const geminiAPIURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key="

// GeminiService encapsulates the logic for interacting with the Google Gemini API.
type GeminiService struct{}

// FactCheck sends a statement to the Gemini API for analysis.
func (s *GeminiService) FactCheck(statement string) (*models.GeminiResponse, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filepath.Dir(filepath.Dir(b)))
	promptPath := filepath.Join(basepath, "prompts", "fact_check_prompt.txt")

	promptBytes, err := os.ReadFile(promptPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read prompt file at %s: %w", promptPath, err)
	}
	promptTemplate := string(promptBytes)
	finalPrompt := fmt.Sprintf(promptTemplate, statement)

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

	return &geminiResp, nil
}
