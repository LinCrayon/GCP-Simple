package image_generation

import (
	"context"
	"fmt"
	"google.golang.org/genai"
	"io"
	"os"
)

/**
 * @Description
 * @Author linshengqian
 * @Date 2026/3/3 16:20
 **/
func generateMMFlashWithGoogleSearch(w io.Writer) error {
	ctx := context.Background()

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Errorf("GOOGLE_API_KEY not set")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})

	if err != nil {
		return fmt.Errorf("创建 genai 客户端失败： %w", err)
	}

	model := "gemini-3.1-flash-image-preview"
	aspectRatio := "16:9" // "1:1","1:4","1:8","2:3","3:2","3:4","4:1","4:3","4:5","5:4","8:1","9:16","16:9","21:9"
	contents := []*genai.Content{
		{
			Role: genai.RoleUser,
			Parts: []*genai.Part{
				{
					Text: "Visualize the current weather forecast for the next 5 days in San Francisco as a clean, modern weather chart. Add a visual on what I should wear each day",
				},
			},
		},
	}

	config := genai.GenerateContentConfig{
		ResponseModalities: []string{
			string(genai.ModalityText),
			string(genai.ModalityImage),
		},
		ImageConfig: &genai.ImageConfig{
			AspectRatio: aspectRatio,
		},
		Tools: []*genai.Tool{
			{
				GoogleSearch: &genai.GoogleSearch{},
			},
		},
		SafetySettings: []*genai.SafetySetting{
			{
				Method:    genai.HarmBlockMethodProbability,
				Category:  genai.HarmCategoryDangerousContent,
				Threshold: genai.HarmBlockThresholdBlockMediumAndAbove,
			},
		},
	}

	resp, err := client.Models.GenerateContent(
		ctx,
		model,
		contents,
		&config,
	)
	if err != nil {
		return fmt.Errorf("无法生成内容： %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return fmt.Errorf("no candidates returned")
	}

	var fileName string
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" {
			fmt.Fprintln(w, part.Text)
		} else if part.InlineData != nil {
			fileName := "weather.png"
			if err := os.WriteFile(fileName, part.InlineData.Data, 0o644); err != nil {
				return fmt.Errorf("保存图像失败：%w", err)
			}
		}
	}
	fmt.Fprintln(w, fileName)

	return nil
}
