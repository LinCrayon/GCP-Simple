package text_generation

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateTextWithPDF PDF文档分析
func generateTextWithPDF(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.5-flash"
	contents := []*genai.Content{
		{Parts: []*genai.Part{
			{Text: `You are a highly skilled document summarization specialist.
	Your task is to provide a concise executive summary of no more than 300 words.
	Please summarize the given document for a general audience.`},
			{FileData: &genai.FileData{
				FileURI:  "gs://cloud-samples-data/generative-ai/pdf/1706.03762v7.pdf",
				MIMEType: "application/pdf",
			}},
		},
			Role: genai.RoleUser},
	}

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, nil)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	respText := resp.Text()

	fmt.Fprintln(w, respText)

	// Example response:
	// "Attention Is All You Need" introduces the Transformer,
	// a groundbreaking neural network architecture designed for...
	// ...

	return nil
}
