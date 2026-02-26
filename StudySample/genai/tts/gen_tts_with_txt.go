package tts

import (
	"context"
	"encoding/binary"
	"google.golang.org/genai"
	"io"
)

func generateTTS(w io.Writer, text string) error {
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
				VoiceConfig: &genai.VoiceConfig{
					PrebuiltVoiceConfig: &genai.PrebuiltVoiceConfig{
						VoiceName: "Kore",
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

func WriteWav(
	w io.Writer,
	pcm []byte,
	channels int,
	sampleRate int,
	bitsPerSample int,
) error {

	byteRate := sampleRate * channels * bitsPerSample / 8
	blockAlign := channels * bitsPerSample / 8
	dataSize := len(pcm)

	// RIFF header
	w.Write([]byte("RIFF"))
	binary.Write(w, binary.LittleEndian, uint32(36+dataSize))
	w.Write([]byte("WAVE"))

	// fmt chunk
	w.Write([]byte("fmt "))
	binary.Write(w, binary.LittleEndian, uint32(16)) // PCM
	binary.Write(w, binary.LittleEndian, uint16(1))  // PCM format
	binary.Write(w, binary.LittleEndian, uint16(channels))
	binary.Write(w, binary.LittleEndian, uint32(sampleRate))
	binary.Write(w, binary.LittleEndian, uint32(byteRate))
	binary.Write(w, binary.LittleEndian, uint16(blockAlign))
	binary.Write(w, binary.LittleEndian, uint16(bitsPerSample))

	// data chunk
	w.Write([]byte("data"))
	binary.Write(w, binary.LittleEndian, uint32(dataSize))
	w.Write(pcm)

	return nil
}
