package image_generation

import (
	"context"
	"fmt"
	"io"
	"os"

	"google.golang.org/genai"
)

// generateMMFlashWithText 以对话的方式生成和处理图像
func generateMMFlashWithText(w io.Writer) error {
	ctx := context.Background()

	//创建连接
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.5-flash-image"
	contents := []*genai.Content{
		{
			Parts: []*genai.Part{
				{Text: "Generate an image of the Eiffel tower with fireworks in the background."},
			},
			Role: genai.RoleUser,
		},
	}

	resp, err := client.Models.GenerateContent(ctx,
		modelName,
		contents,
		&genai.GenerateContentConfig{
			ResponseModalities: []string{ //返回内容类型
				string(genai.ModalityText),
				string(genai.ModalityImage),
			},
			CandidateCount: int32(1), //模型一次生成几个候选结果
			SafetySettings: []*genai.SafetySetting{
				{Method: genai.HarmBlockMethodProbability},
				{Category: genai.HarmCategoryDangerousContent},
				{Threshold: genai.HarmBlockThresholdBlockMediumAndAbove},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return fmt.Errorf("no candidates returned")
	}
	var fileName string
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" {
			fmt.Fprintln(w, part.Text)
		} else if part.InlineData != nil {
			fileName = "example-image-eiffel-tower.png"
			if err := os.WriteFile(fileName, part.InlineData.Data, 0o644); err != nil {
				return fmt.Errorf("failed to save image: %w", err)
			}
		}
	}
	fmt.Fprintln(w, fileName)

	// 响应示例：我将生成一张埃菲尔铁塔在夜间的图像，其背后的黑暗天空中绽放着绚丽多彩的烟花。
	return nil
}
