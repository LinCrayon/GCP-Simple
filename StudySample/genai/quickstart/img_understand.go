package quickstart

import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateWithTextImage 图像理解
func generateWithTextImage(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("创建 genai 客户端失败: %w", err)
	}

	modelName := "gemini-2.5-flash"
	contents := []*genai.Content{
		{
			Parts: []*genai.Part{
				{Text: "这张图片显示了什么?"},
				{FileData: &genai.FileData{
					// Image source: https://storage.googleapis.com/cloud-samples-data/generative-ai/image/scones.jpg
					FileURI:  "gs://training-lsq-week3-iam-visible/1.png",
					MIMEType: "image/png",
				}},
			},
			Role: genai.RoleUser,
		},
	}

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, &genai.GenerateContentConfig{
		ResponseModalities: []string{string(genai.ModalityText)},
		CandidateCount:     1,
	})
	if err != nil {
		return fmt.Errorf("生成内容失败: %w", err)
	}

	// 遍历候选结果，输出文本
	for _, candidate := range resp.Candidates {
		if candidate.Content == nil {
			continue
		}
		for _, part := range candidate.Content.Parts {
			if part.Text != "" {
				fmt.Fprintln(w, part.Text)
			}
		}
	}
	return nil
}
