package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const apiURL = "https://api.openai.com/v1/chat/completions"

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"` // Could be used later for truncation check
	} `json:"choices"`
}

const maxInputChars = 3000 // ~750 tokens roughly

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("🔑 Enter your OpenAI API key: ")
	apiKeyInput, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("❌ Failed to read API key:", err)
		os.Exit(1)
	}
	apiKey := strings.TrimSpace(apiKeyInput)
	if apiKey == "" {
		fmt.Println("❌ API key is required.")
		os.Exit(1)
	}

	fmt.Println("🤖 ChatGPT REPL (type 'exit' to quit, 'clear' to reset)")

	history := []ChatMessage{
		{Role: "system", Content: "You are a helpful assistant."},
	}

	for {
		fmt.Print(">> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "exit", "quit":
			fmt.Println("👋 Goodbye!")
			return
		case "clear":
			// Reset history except system prompt
			history = history[:1]
			fmt.Println("🧹 Conversation history cleared.")
			continue
		}

		// Check input length and warn if too long
		if len(input) > maxInputChars {
			fmt.Printf("⚠️ Warning: Input is very long (%d chars). Consider shortening.\n", len(input))
		}

		history = append(history, ChatMessage{Role: "user", Content: input})

		requestBody := ChatRequest{
			Model:    "gpt-3.5-turbo",
			Messages: history,
		}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			fmt.Println("❌ Failed to marshal request:", err)
			continue
		}

		req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("❌ Failed to create request:", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("❌ Request failed:", err)
			continue
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("❌ API error: %s\n%s\n", resp.Status, string(body))
			continue
		}

		var chatResp ChatResponse
		if err := json.Unmarshal(body, &chatResp); err != nil {
			fmt.Println("❌ Failed to decode response:", err)
			continue
		}

		if len(chatResp.Choices) == 0 {
			fmt.Println("❌ No choices returned from API.")
			continue
		}

		reply := chatResp.Choices[0].Message.Content
		fmt.Println("GPT:", reply)
		history = append(history, ChatMessage{
			Role:    chatResp.Choices[0].Message.Role,
			Content: reply,
		})
	}
}
