package live

// [START googlegenaisdk_live_txt_with_audio]
import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"google.golang.org/genai"
)

// generateLiveTextWithAudio 音频输入，文本输出
// 发送音频，接收文本回复
func generateLiveTextWithAudio(w io.Writer) error {
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
	}

	session, err := client.Live.Connect(ctx, modelName, config)
	if err != nil {
		return fmt.Errorf("failed to connect live: %w", err)
	}
	defer session.Close()

	audioURL := "https://storage.googleapis.com/generativeai-downloads/data/16000.wav"
	// Download audio
	resp, err := http.Get(audioURL)
	if err != nil {
		return fmt.Errorf("failed to download audio: %w", err)
	}
	defer resp.Body.Close()

	audioBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read audio: %w", err)
	}

	fmt.Fprintf(w, "> Answer to this audio url: %s\n\n", audioURL)

	// Send the audio as Blob media input
	err = session.SendRealtimeInput(genai.LiveRealtimeInput{
		Media: &genai.Blob{
			Data:     audioBytes,
			MIMEType: "audio/pcm;rate=16000",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send audio input: %w", err)
	}

	// Stream the response
	var response strings.Builder
	for {
		chunk, err := session.Receive()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error receiving response: %w", err)
		}

		if chunk.ServerContent == nil {
			continue
		}

		// Handle model turn responses
		if chunk.ServerContent.ModelTurn != nil {
			for _, part := range chunk.ServerContent.ModelTurn.Parts {
				if part != nil && part.Text != "" {
					response.WriteString(part.Text)
				}
			}
		}
	}

	fmt.Fprintln(w, response.String())

	// Example output:
	// > Answer to this audio url: https://storage.googleapis.com/generativeai-downloads/data/16000.wav
	// Yes, I can hear you. How can I help you today?
	return nil
}
