package text_generation

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateChatWithText 多轮对话：保持上下文的对话、记忆历史消息
func generateChatWithText(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}
	modelName := "gemini-2.5-flash"
	history := []*genai.Content{
		{
			Role: genai.RoleUser,
			Parts: []*genai.Part{
				{Text: "Hello there"},
			},
		},
		{
			Role: "model",
			Parts: []*genai.Part{
				{Text: "Great to meet you. What would you like to know?"},
			},
		},
	}
	chatSession, err := client.Chats.Create(ctx, modelName, nil, history)
	if err != nil {
		return fmt.Errorf("failed to create genai chat session: %w", err)
	}
	contents := genai.Part{Text: "Tell me a story."}
	resp, err := chatSession.SendMessage(ctx, contents)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	respText := resp.Text()

	fmt.Fprintln(w, respText)
	// Example response:
	// Okay, settle in. Let me tell you a story about a quiet cartographer, but not of lands and seas.
	// ...
	// In the sleepy town of Oakhaven, nestled between the Whispering Hills and the Murmuring River, lived a woman named Elara.
	// ...

	return nil
}
