package text_generation

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateWithRouting 模型路由：自动选择最佳模型
func generateWithRouting(w io.Writer) error {
	ctx := context.Background()

	clientConfig := &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1beta1"},
	}

	client, err := genai.NewClient(ctx, clientConfig)

	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelSelectionConfig := &genai.ModelSelectionConfig{
		FeatureSelectionPreference: genai.FeatureSelectionPreferencePrioritizeQuality,
	}

	generateContentConfig := &genai.GenerateContentConfig{
		ModelSelectionConfig: modelSelectionConfig,
	}

	modelName := "model-optimizer-exp-04-09"

	resp, err := client.Models.GenerateContent(ctx,
		modelName,
		genai.Text("How does AI work?"),
		generateContentConfig,
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
