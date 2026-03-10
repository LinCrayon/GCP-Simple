package live

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateLiveCodeExecWithTxt 代码执行：让模型执行代码并返回结果
// that performs code exec with text calls and handles responses.
func generateLiveCodeExecWithTxt(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.0-flash-live-preview-04-09"

	config := &genai.LiveConnectConfig{
		ResponseModalities: []genai.Modality{genai.ModalityText},
		Tools: []*genai.Tool{
			{
				CodeExecution: &genai.ToolCodeExecution{},
			},
		},
	}

	session, err := client.Live.Connect(ctx, modelName, config)
	if err != nil {
		return fmt.Errorf("failed to connect live session: %w", err)
	}
	defer session.Close()

	textInput := "Compute the largest prime palindrome under 10"
	fmt.Fprintf(w, "> %s\n\n", textInput)

	err = session.SendClientContent(genai.LiveClientContentInput{
		Turns: []*genai.Content{
			{
				Role: genai.RoleUser,
				Parts: []*genai.Part{
					{Text: textInput},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send client content: %w", err)
	}

	var response string

	// Receive streaming responses
	for {
		chunk, err := session.Receive()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error receiving stream: %w", err)
		}

		// Handle the main model output
		if chunk.ServerContent != nil {
			if chunk.ServerContent.ModelTurn != nil {
				for _, part := range chunk.ServerContent.ModelTurn.Parts {
					if part == nil {
						continue
					}
					if part.ExecutableCode != nil {
						fmt.Fprint(w, part.ExecutableCode.Code)
					}
					if part.CodeExecutionResult != nil {
						response += part.CodeExecutionResult.Output
					}
				}

			}
		}
	}
	// Example output:
	//  > Compute the largest prime palindrome under 10
	//  Final Answer: The final answer is $\boxed{7}$
	fmt.Fprintln(w, response)
	return nil
}
