package tts

import (
	"bytes"
	"os"
	"testing"
)

func TestVideoGeneration(t *testing.T) {

	t.Run("TextToSpeech", func(t *testing.T) {
		var pcm bytes.Buffer

		err := generateTTS(&pcm, "Have a wonderful day!")
		if err != nil {
			t.Fatalf("generateTTS failed: %v", err)
		}

		if pcm.Len() == 0 {
			t.Fatal("expected audio output, got empty")
		}

		_ = os.MkdirAll("testdata", 0755)

		f, err := os.Create("testdata/out.wav")
		if err != nil {
			t.Fatalf("create wav failed: %v", err)
		}

		defer f.Close()

		err = WriteWav(
			f,
			pcm.Bytes(),
			1,     // mono
			24000, // sample rate (Gemini TTS 默认)
			16,    // bits per sample
		)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("MultiSpeakerTTS", func(t *testing.T) {
		prompt := `Make Speaker1 sound tired and bored, and Speaker2 sound excited and happy:

Speaker1: So... what's on the agenda today?
Speaker2: You're never going to guess!?`

		err := os.MkdirAll("testdata", 0755)
		if err != nil {
			t.Fatal(err)
		}

		var pcm bytes.Buffer

		err = generateMultiSpeakerTTS(&pcm, prompt)
		if err != nil {
			t.Fatalf("generateMultiSpeakerTTS failed: %v", err)
		}

		if pcm.Len() == 0 {
			t.Fatal("expected audio output, got empty")
		}

		f, err := os.Create("testdata/multi_speaker.wav")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		err = WriteWav(
			f,
			pcm.Bytes(),
			1,     // mono
			24000, // Gemini TTS 默认采样率
			16,    // PCM16
		)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("multi-speaker TTS written to testdata/multi_speaker.wav")
	})

	t.Run("生成播客文本转为tts", func(t *testing.T) {
		transcript := `
Dr. Anya: This frog species is absolutely incredible!
Liam: I know! The coloration alone is enough to blow your mind.
Dr. Anya: And the call patterns are unlike anything we’ve recorded before.
`

		var pcm bytes.Buffer

		err := GenerateMultiSpeakerTTS(transcript, &pcm)
		if err != nil {
			t.Fatalf("GenerateMultiSpeakerTTS failed: %v", err)
		}

		if pcm.Len() == 0 {
			t.Fatal("expected audio output, got empty buffer")
		}

		// 确保 testdata 目录存在
		_ = os.MkdirAll("testdata", 0755)

		f, err := os.Create("testdata/multi_speaker.wav")
		if err != nil {
			t.Fatalf("create wav failed: %v", err)
		}
		defer f.Close()

		err = WriteWav(
			f,
			pcm.Bytes(),
			1,     // mono
			24000, // Gemini TTS 默认
			16,    // PCM16
		)
		if err != nil {
			t.Fatalf("WriteWav failed: %v", err)
		}
	})
}
