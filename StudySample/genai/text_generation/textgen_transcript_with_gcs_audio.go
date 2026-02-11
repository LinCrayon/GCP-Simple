package text_generation

import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateAudioTranscript 音频转录：语音转文字 （字幕生成）
func generateAudioTranscript(w io.Writer) error {
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
			{Text: `Transcribe the interview, in the format of timecode, speaker, caption.
Use speaker A, speaker B, etc. to identify speakers.`},
			{FileData: &genai.FileData{
				FileURI:  "gs://cloud-samples-data/generative-ai/audio/pixel.mp3",
				MIMEType: "audio/mpeg",
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
	// 00:00:00, A: your devices are getting better over time.
	// 00:01:13, A: And so we think about it across the entire portfolio from phones to watch, ...
	// ...

	return nil
}
