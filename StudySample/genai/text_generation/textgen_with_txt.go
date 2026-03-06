package text_generation

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateWithText 简单的文本生成
func generateWithText(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	val := int32(2000)
	resp, err := client.Models.GenerateContent(ctx,
		"gemini-2.5-flash",
		genai.Text("How does AI work?"),
		&genai.GenerateContentConfig{
			ThinkingConfig: &genai.ThinkingConfig{
				IncludeThoughts: true,                           //是否包含思考过程,答案包含(逻辑推导、最终答案)
				ThinkingBudget:  &val,                           //思考预算/Token 限制
				ThinkingLevel:   genai.ThinkingLevelUnspecified, //思考等级/深度 (自动选择根据模型选择平衡点)
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	respText := resp.Text()

	fmt.Fprintln(w, respText)
	// Example response:
	// That's a great question! Understanding how AI works can feel like ...
	// ...
	// **1. The Foundation: Data and Algorithms**
	// ...

	return nil
}
