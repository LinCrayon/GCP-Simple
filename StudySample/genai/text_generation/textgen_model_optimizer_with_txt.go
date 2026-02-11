package text_generation

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateModelOptimizerWithTxt shows how to generate text using a text prompt and model optimizer.
func generateModelOptimizerWithTxt(w io.Writer) error {
	ctx := context.Background()

	clientConfig := &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1beta1"},
	}

	client, err := genai.NewClient(ctx, clientConfig)

	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelSelectionConfig := &genai.ModelSelectionConfig{
		FeatureSelectionPreference: genai.FeatureSelectionPreferenceBalanced,
	}

	generateContentConfig := &genai.GenerateContentConfig{
		ModelSelectionConfig: modelSelectionConfig,
	}

	modelName := "gemini-2.5-flash"
	contents := genai.Text("How does AI work?")

	resp, err := client.Models.GenerateContent(ctx,
		modelName,
		contents,
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
