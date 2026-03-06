package text_generation

import (
	"context"
	"fmt"
	"google.golang.org/genai"
	"io"
)

/**
 * @Description
 * @Author linshengqian
 * @Date 2026/3/6 10:06
 **/
func generateWithRespSchema(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: "array",
			Items: &genai.Schema{
				Type: "object",
				Properties: map[string]*genai.Schema{
					"recipe_name": {Type: "string"},
					"ingredients": {
						Type:  "array",
						Items: &genai.Schema{Type: "string"},
					},
				},
				Required: []string{"recipe_name", "ingredients"},
			},
		},
	}
	/*
		config := &genai.GenerateContentConfig{
				ResponseMIMEType: "application/json",
				ResponseSchema: &genai.Schema{
					Type: "object",
					Properties: map[string]*genai.Schema{
						"forecast": {
							Type: "array",
							Items: &genai.Schema{
								Type: "object",
								Properties: map[string]*genai.Schema{
									"Day":         {Type: "string", Nullable: genai.Ptr(true)},
									"Forecast":    {Type: "string", Nullable: genai.Ptr(true)},
									"Temperature": {Type: "integer", Nullable: genai.Ptr(true)},
									"Humidity":    {Type: "string", Nullable: genai.Ptr(true)},
									"Wind Speed":  {Type: "integer", Nullable: genai.Ptr(true)},
								},
								Required: []string{"Day", "Temperature", "Forecast", "Wind Speed"},
							},
						},
					},
				},
			}
	*/

	/*
		config := &genai.GenerateContentConfig{
			ResponseMIMEType: "text/x.enum",
			ResponseSchema: &genai.Schema{
				Type: "STRING",
				Enum: []string{"Percussion", "String", "Woodwind", "Brass", "Keyboard"},
			},
		}
	*/

	contents := []*genai.Content{
		{Parts: []*genai.Part{
			{Text: "List a few popular cookie recipes."},
		},
			Role: "user"},
	}
	modelName := "gemini-2.5-flash"

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	respText := resp.Text()

	fmt.Fprintln(w, respText)

	return nil
}
