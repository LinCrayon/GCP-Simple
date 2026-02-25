package text_generation

import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateWithTextImage 多模态输入:文本 + 图片 (问："图片中有什么？" + 一张图片)
func generateWithTextImage(w io.Writer) error {
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
			{Text: "这张图片显示了什么？"},
			{FileData: &genai.FileData{
				// Image source: https://storage.googleapis.com/cloud-samples-data/generative-ai/image/scones.jpg
				FileURI:  "gs://cloud-samples-data/generative-ai/image/scones.jpg",
				MIMEType: "image/jpeg",
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
	// The image shows an overhead shot of a rustic, artistic arrangement on a surface that ...

	return nil
}
