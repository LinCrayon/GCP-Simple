package text_generation

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateChatStreamWithText shows how to generate chat stream using a text prompt.
func generateChatStreamWithText(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.5-flash"

	chatSession, err := client.Chats.Create(ctx, modelName, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create genai chat session: %w", err)
	}

	var streamErr error
	contents := genai.Part{Text: "Why is the sky blue?"}

	stream := chatSession.SendMessageStream(ctx, contents)
	stream(func(resp *genai.GenerateContentResponse, err error) bool {
		if err != nil {
			streamErr = err
			return false
		}
		for _, cand := range resp.Candidates {
			for _, part := range cand.Content.Parts {
				fmt.Fprintln(w, part.Text)
			}
		}
		return true
	})

	// Example response:
	// The
	// sky appears blue due to a phenomenon called **Rayleigh scattering**.
	// Here's a breakdown:
	// ...

	return streamErr
}
