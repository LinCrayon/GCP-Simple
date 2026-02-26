package tts

import (
	"context"
	"io"

	"google.golang.org/genai"
)

func GenerateMultiSpeakerTTS(text string, w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return err
	}

	resp, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-preview-tts",
		genai.Text(text),
		&genai.GenerateContentConfig{
			ResponseModalities: []string{"AUDIO"},
			SpeechConfig: &genai.SpeechConfig{
				MultiSpeakerVoiceConfig: &genai.MultiSpeakerVoiceConfig{
					SpeakerVoiceConfigs: []*genai.SpeakerVoiceConfig{
						{
							Speaker: "Dr. Anya",
							VoiceConfig: &genai.VoiceConfig{
								PrebuiltVoiceConfig: &genai.PrebuiltVoiceConfig{
									VoiceName: "Kore",
								},
							},
						},
						{
							Speaker: "Liam",
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
		return err
	}

	audio := resp.Candidates[0].
		Content.Parts[0].
		InlineData.Data

	_, err = w.Write(audio)
	return err
}
