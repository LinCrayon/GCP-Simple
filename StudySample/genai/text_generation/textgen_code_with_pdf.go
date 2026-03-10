package text_generation

import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateWithPDF 代码转换： 输入PDF中的代码，输出符合规范的代码
func generateWithPDF(w io.Writer) error {
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
			{Text: "Convert this python code to use Google Python Style Guide."},
			{FileData: &genai.FileData{
				FileURI:  "https://storage.googleapis.com/cloud-samples-data/generative-ai/text/inefficient_fibonacci_series_python_code.pdf",
				MIMEType: "application/pdf",
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
	// ```python
	// def fibonacci(n: int) -> list[int]:
	// ...

	return nil
}
