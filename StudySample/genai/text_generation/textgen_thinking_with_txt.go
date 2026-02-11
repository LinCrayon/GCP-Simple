package text_generation

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateThinkingWithText 思维过程：展示推理过程
func generateThinkingWithText(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	resp, err := client.Models.GenerateContent(ctx,
		"gemini-2.5-flash",
		genai.Text("solve x^2 + 4x + 4 = 0"),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	respText := resp.Text()

	fmt.Fprintln(w, respText)
	// Example response:
	// To solve the quadratic equation $x^2 + 4x + 4 = 0$, we can use a few methods:
	//
	// **Method 1: Factoring (Recognizing a Perfect Square Trinomial)**
	// **1. The Foundation: Data and Algorithms**
	//
	// Notice that the left side of the equation is a perfect square trinomial.
	// ...

	return nil
}
