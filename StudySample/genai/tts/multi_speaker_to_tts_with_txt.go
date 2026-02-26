package tts

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

func generateMultiSpeakerTTS(w io.Writer, prompt string) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return err
	}

	resp, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-preview-tts",
		genai.Text(prompt),
		&genai.GenerateContentConfig{
			ResponseModalities: []string{"AUDIO"},
			SpeechConfig: &genai.SpeechConfig{
				MultiSpeakerVoiceConfig: &genai.MultiSpeakerVoiceConfig{
					SpeakerVoiceConfigs: []*genai.SpeakerVoiceConfig{
						{
							Speaker: "Joe",
							VoiceConfig: &genai.VoiceConfig{
								PrebuiltVoiceConfig: &genai.PrebuiltVoiceConfig{
									VoiceName: "Kore",
								},
							},
						},
						{
							Speaker: "Jane",
							VoiceConfig: &genai.VoiceConfig{
								PrebuiltVoiceConfig: &genai.PrebuiltVoiceConfig{
									VoiceName: "Puck",
								},
							},
						},
					},
				},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("GenerateContent failed: %w", err)
	}

	audio := resp.
		Candidates[0].
		Content.Parts[0].
		InlineData.Data

	_, err = w.Write(audio)
	return err
}
