package models

// NEW: Source defines the structure for a single search result source.
// This will be sent to the frontend for display.
type Source struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// FactCheckRequest defines the structure for incoming API requests.
type FactCheckRequest struct {
	Statement string `json:"statement"`
}

// GeminiResponse defines the structure of the JSON response sent back to the client.
// It now includes the list of sources used for the fact-check.
type GeminiResponse struct {
	Verdict           string   `json:"verdict"`
	Confidence        string   `json:"confidence"`
	Reason            string   `json:"reason"`
	AdditionalContext string   `json:"additional_context"`
	Sources           []Source `json:"sources"` // The new field for sources
}

// FactCheckHistoryItem defines the structure for a single record from the database.
type FactCheckHistoryItem struct {
	ID                int    `json:"id"`
	Statement         string `json:"statement"`
	Verdict           string `json:"verdict"`
	Confidence        string `json:"confidence"`
	Reason            string `json:"reason"`
	AdditionalContext string `json:"additional_context"`
	CreatedAt         string `json:"created_at"`
}
