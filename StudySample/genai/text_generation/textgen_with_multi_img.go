package text_generation

import (
	"context"
	"fmt"
	"io"
	"os"

	genai "google.golang.org/genai"
)

// generateWithMultiImg 文本+多张图片
func generateWithMultiImg(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	// TODO(Developer): Update the path to file (image source:
	//   https://storage.googleapis.com/cloud-samples-data/generative-ai/image/latte.jpg )
	imageBytes, err := os.ReadFile("./latte.jpg")
	if err != nil {
		return fmt.Errorf("failed to read image: %w", err)
	}

	contents := []*genai.Content{
		{Parts: []*genai.Part{
			{Text: "Write an advertising jingle based on the items in both images."},
			{FileData: &genai.FileData{
				// Image source: https://storage.googleapis.com/cloud-samples-data/generative-ai/image/scones.jpg
				FileURI:  "gs://cloud-samples-data/generative-ai/image/scones.jpg",
				MIMEType: "image/jpeg",
			}},
			{InlineData: &genai.Blob{
				Data:     imageBytes,
				MIMEType: "image/jpeg",
			}},
		}},
	}
	modelName := "gemini-2.5-flash"

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, nil)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	respText := resp.Text()

	fmt.Fprintln(w, respText)

	// Example response:
	// Okay, here's an advertising jingle inspired by the blueberry scones, coffee, flowers, chocolate cake, and latte:
	//
	// (Upbeat, jazzy music)
	// ...

	return nil
}
