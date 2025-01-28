package database

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"gorm.io/gorm"

	openai "github.com/sashabaranov/go-openai"
)

// Global variable for the GPT client.
// Because this example is for local/small projects, we're not using concurrency locks.
var gptClient *openai.Client

// ---------------------------------------------------------------------
// 1. InitializeGPT - sets up the client once at startup
// ---------------------------------------------------------------------
func InitializeGPT(apiKey, baseURL string) error {
	if apiKey == "fill back later" {
		// Fall back to environment variable if no apiKey was passed in
		apiKey = os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return fmt.Errorf("API key is not provided or set in environment")
		}
	}

	// Default configuration from the library
	config := openai.DefaultConfig(apiKey)

	// If you need a custom base URL (e.g., "https://api.zhizengzeng.com/v1"), set it here
	if baseURL != "" {
		config.BaseURL = baseURL
	} else {
		config.BaseURL = "https://api.zhizengzeng.com/v1"
	}

	// Create the global client
	gptClient = openai.NewClientWithConfig(config)
	return nil
}

// ---------------------------------------------------------------------
// 2. QueryGPT - sends a user prompt and returns the AI's response
// ---------------------------------------------------------------------
func QueryGPT(prompt string) (string, error) {
	if gptClient == nil {
		return "", fmt.Errorf("GPT client is not initialized")
	}

	// Create the request
	req := openai.ChatCompletionRequest{
		Model: "gpt-o1", // or "SparkDesk" or whatever your provider expects
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a helpful assistant.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		// Optional: Temperature, MaxTokens, etc.
	}

	// Perform the call
	resp, err := gptClient.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("CreateChatCompletion error: %v", err)
	}

	// Extract the AI's text from the first choice
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned in response")
	}
	return resp.Choices[0].Message.Content, nil
}

// ---------------------------------------------------------------------
// 3. (Optional) GPTHandler - an HTTP handler for receiving queries
// ---------------------------------------------------------------------
func PostGPTHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := CheckCookie(r.Header.Get("Cookie"), db); err != nil {
		http.Error(w, "invalid cookie", http.StatusUnauthorized)
		return
	}

	// Define a small struct to parse the incoming JSON
	type requestBody struct {
		Prompt string `json:"prompt"`
	}

	var body requestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil || body.Prompt == "" {
		http.Error(w, "Invalid JSON or missing 'prompt'", http.StatusBadRequest)
		return
	}

	// Call our QueryGPT function
	answer, err := QueryGPT(body.Prompt)
	if err != nil {
		http.Error(w, "GPT query error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with JSON
	respData := map[string]string{"answer": answer}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(respData)
}
