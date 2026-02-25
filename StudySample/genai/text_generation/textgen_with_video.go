package text_generation

import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateWithVideo 视频内容分析 输出：总结 + 章节时间戳
func generateWithVideo(w io.Writer) error {
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
			{Text: `分析提供的视频文件，包括其音频。
简要概括视频的要点。
使用时间戳为所讨论的关键部分或主题创建章节细分。`},
			{FileData: &genai.FileData{
				FileURI:  "gs://cloud-samples-data/generative-ai/video/pixel8.mp4",
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
	// Here's an analysis of the provided video file:
	//
	// **Summary**
	//
	// The video features Saeka Shimada, a photographer in Tokyo, who uses the new Pixel phone ...
	//
	// **Chapter Breakdown**
	//
	// *   **0:00-0:05**: Introduction to Saeka Shimada and her work as a photographer in Tokyo.
	// ...

	return nil
}

// [END googlegenaisdk_textgen_with_video]
