package text_generation

import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateWithYTVideo 展示如何使用 YouTube 视频作为输入生成文本。
func generateWithYTVideo(w io.Writer) error {
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
			{Text: "根据该视频撰写一篇简短且引人入胜的博客文章。"},
			{FileData: &genai.FileData{
				FileURI:  "https://www.youtube.com/watch?v=3KtWfp0UopM",
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
	// Okay, here’s a short and engaging blog post based on the provided video.
	//
	// **Google's 25th: A Look Back at What We've Searched**
	// ...

	return nil
}
