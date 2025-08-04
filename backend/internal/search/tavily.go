package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const tavilyAPIURL = "https://api.tavily.com/search"

// TavilySearchResult defines the structure for a single search result from the API.
type TavilySearchResult struct {
	Title   string  `json:"title"`
	URL     string  `json:"url"`
	Content string  `json:"content"`
	Score   float64 `json:"score"`
}

// TavilyResponse defines the overall structure of the response from the Tavily API.
type TavilyResponse struct {
	Query   string               `json:"query"`
	Results []TavilySearchResult `json:"results"`
}

// TavilyService encapsulates the logic for interacting with the Tavily Search API.
type TavilyService struct{}

// Search performs a web search using the Tavily API.
// It takes a query string and returns a slice of search results.
func (s *TavilyService) Search(query string) ([]TavilySearchResult, error) {
	// 1. Load the Tavily API key from environment variables.
	apiKey := os.Getenv("TAVILY_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("TAVILY_API_KEY environment variable not set")
	}

	// 2. Construct the request body for the Tavily API.
	// We configure it to get 3 results and include the content of the pages.
	requestPayload := map[string]interface{}{
		"api_key":      apiKey,
		"query":        query,
		"search_depth": "basic",
		"max_results":  3, // We'll get the top 3 most relevant results.
	}

	requestBody, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tavily request body: %w", err)
	}

	// 3. Create and send the HTTP POST request with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", tavilyAPIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create tavily http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to tavily api: %w", err)
	}
	defer resp.Body.Close()

	// 4. Read and parse the response.
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read tavily response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tavily api returned non-200 status: %s", string(responseBody))
	}

	var tavilyResp TavilyResponse
	if err := json.Unmarshal(responseBody, &tavilyResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tavily response: %w", err)
	}

	// 5. Return the slice of results.
	return tavilyResp.Results, nil
}
