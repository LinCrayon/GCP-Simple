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
			{Text: `您是一位技术精湛的文档摘要专家。
	您的任务是提供不超过 300 字的简洁执行摘要。
	请为一般读者总结给定的文件。`},
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
