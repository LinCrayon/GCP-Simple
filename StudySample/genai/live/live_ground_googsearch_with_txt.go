package live

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateGroundSearchWithTxt 谷歌搜索：在回答中引用谷歌搜索结果
func generateGroundSearchWithTxt(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.0-flash-live-preview-04-09"

	config := &genai.LiveConnectConfig{
		ResponseModalities: []genai.Modality{genai.ModalityText},
		Tools: []*genai.Tool{
			{GoogleSearch: &genai.GoogleSearch{}},
		},
	}

	session, err := client.Live.Connect(ctx, modelName, config)
	if err != nil {
		return fmt.Errorf("failed to connect live session: %w", err)
	}
	defer session.Close()

	textInput := "When did the last Brazil vs. Argentina soccer match happen?"

	// Send user input
	userContent := &genai.Content{
		Role: genai.RoleUser,
		Parts: []*genai.Part{
			{Text: textInput},
		},
	}
	if err := session.SendClientContent(genai.LiveClientContentInput{
		Turns: []*genai.Content{userContent},
	}); err != nil {
		return fmt.Errorf("failed to send client content: %w", err)
	}

	var response string

	// Receive streaming responses
	for {
		chunk, err := session.Receive()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error receiving stream: %w", err)
		}

		// Handle the main model output
		if chunk.ServerContent != nil {
			if chunk.ServerContent.ModelTurn != nil {
				for _, part := range chunk.ServerContent.ModelTurn.Parts {
					if part == nil {
						continue
					}
					if part.Text != "" {
						response += part.Text
					}
				}
			}
		}

		if chunk.GoAway != nil {
			break
		}
	}

	fmt.Fprintln(w, response)

	// Example output:
	// > When did the last Brazil vs. Argentina soccer match happen?
	// The most recent match between Argentina and Brazil took place on March 25, 2025, as part of the 2026 World Cup qualifiers. Argentina won 4-1.

	return nil
}
