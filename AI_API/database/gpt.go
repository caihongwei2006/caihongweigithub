package database

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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

	// Default configuration from the library
	config := openai.DefaultConfig(apiKey)

	// If you need a custom base URL (e.g., "https://api.zhizengzeng.com/v1"), set it here
	if baseURL != "" {
		config.BaseURL = baseURL
	} else {
		config.BaseURL = "https://api.zhizengzeng.com/v1" // default base URL
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

	// 创建请求，仅包含 user 消息，并设置受限参数为固定值
	req := openai.ChatCompletionRequest{
		Model: "o1-mini", // 确认使用正确的模型名称
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		// 设置受限参数为固定值
		Temperature:         1.0,
		TopP:                1.0,
		N:                   1,
		PresencePenalty:     0.0,
		FrequencyPenalty:    0.0,
		MaxCompletionTokens: 15000, // 可以保留或根据需要调整
	}

	// 发起请求
	resp, err := gptClient.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("CreateChatCompletion error: %v", err)
	}

	// 提取第一个选择
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
