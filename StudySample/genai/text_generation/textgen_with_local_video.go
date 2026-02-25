package text_generation

import (
	"context"
	"fmt"
	"io"
	"os"

	genai "google.golang.org/genai"
)

// generateWithLocalVideo 展示如何使用本地视频输入生成文本。
func generateWithLocalVideo(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	// Read local video file content
	data, err := os.ReadFile("describe_video_content.mp4")
	if err != nil {
		return fmt.Errorf("failed to read local video: %w", err)
	}

	modelName := "gemini-2.5-flash"
	contents := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{Text: `根据该视频撰写一篇简短且引人入胜的博客文章。`},
				{InlineData: &genai.Blob{
					MIMEType: "video/mp4",
					Data:     data,
				}},
			},
		},
	}

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, nil)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	respText := resp.Text()
	fmt.Fprintln(w, respText)

	// Example response:
	// Finding Your Flow: The Focused Ascent
	//
	// Ever watched someone scale an indoor climbing wall and been captivated by their precision and power? This video perfectly captures that intense focus and calculated movement.
	//
	// Our climber isn't just pulling himself up; he's engaging in a dynamic dance with gravity. Every reach, every foot placement, every clip of the rope is a deliberate part of solving the route's puzzle. You can almost feel the concentration as his eyes scan for the next optimal hold, his muscles working in unison to propel him upwards.
	//
	// Indoor climbing....
	// ...

	return nil
}
