package text_generation

import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateWithMuteVideo shows how to generate text using a video with no sound as the input.
func generateWithMuteVideo(w io.Writer) error {
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
			{Text: "What is in the video?"},
			{FileData: &genai.FileData{
				FileURI:  "gs://cloud-samples-data/generative-ai/video/ad_copy_from_video.mp4",
				MIMEType: "video/mp4",
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
	// The video shows several surfers riding waves in an ocean setting. The waves are ...

	return nil
}
