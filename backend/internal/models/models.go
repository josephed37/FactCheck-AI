package models

// FactCheckRequest defines the structure for incoming API requests.
// This is the data we expect from the client (our Streamlit app).
type FactCheckRequest struct {
	Statement string `json:"statement"`
}

// GeminiResponse defines the structure of the JSON response we expect from Gemini.
// This is our "data contract" with the AI, similar to our Pydantic model.
type GeminiResponse struct {
	Verdict           string `json:"verdict"`
	Confidence        string `json:"confidence"`
	Reason            string `json:"reason"`
	AdditionalContext string `json:"additional_context"`
}

// retrieved from our database's history. It includes fields that the
// database generates, like ID and CreatedAt.
type FactCheckHistoryItem struct {
	ID                int    `json:"id"`
	Statement         string `json:"statement"`
	Verdict           string `json:"verdict"`
	Confidence        string `json:"confidence"`
	Reason            string `json:"reason"`
	AdditionalContext string `json:"additional_context"`
	CreatedAt         string `json:"created_at"`
}
