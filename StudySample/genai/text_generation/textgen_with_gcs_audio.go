package text_generation

import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateWithAudio 音频分析：音频内容分析
func generateWithAudio(w io.Writer) error {
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
			{Text: `提供音频文件的摘要。
简要概括音频的要点。
使用时间戳为所讨论的关键部分或主题创建章节细分。`},
			{FileData: &genai.FileData{
				FileURI:  "gs://cloud-samples-data/generative-ai/audio/pixel.mp3",
				MIMEType: "audio/mpeg",
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
	// Here is a summary and chapter breakdown of the audio file:
	//
	// **Summary:**
	//
	// The audio file is a "Made by Google" podcast episode discussing the Pixel Feature Drops, ...
	//
	// **Chapter Breakdown:**
	//
	// *   **0:00 - 0:54:** Introduction to the podcast and guests, Aisha Sharif and DeCarlos Love.
	// ...

	return nil
}
