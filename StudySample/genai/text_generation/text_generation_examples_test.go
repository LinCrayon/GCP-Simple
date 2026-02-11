package text_generation

import (
	"bytes"
	"testing"
)

func TestGenerate(t *testing.T) {
	//t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	buf := new(bytes.Buffer)

	t.Run("文本生成-流式回答 Why is the sky blue", func(t *testing.T) {
		buf.Reset()

		err := generateWithTextStream(buf)
		if err != nil {
			t.Fatalf("生成文本流失败: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Fatal("expected some text output, got empty")
		}

		t.Logf("模型输出:\n%s", output)
	})

	t.Run("使用自定义配置生成文本内容", func(t *testing.T) {

		err := generateWithConfig(buf)
		if err != nil {
			t.Fatalf("文本生成失败: %v", err)
		}

		if buf.String() == "" {
			t.Fatal("期望有输出内容，但结果为空")
		}
		t.Logf("模型输出:\n%s", buf.String())
	})

	t.Run("generate with text prompt", func(t *testing.T) {
		buf.Reset()
		err := generateWithText(buf)
		if err != nil {
			t.Fatalf("generateWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate with text prompt and system instructions", func(t *testing.T) {
		buf.Reset()
		err := generateWithSystem(buf)
		if err != nil {
			t.Fatalf("generateWithSystem failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate with text and image prompt", func(t *testing.T) {
		buf.Reset()
		err := generateWithTextImage(buf)
		if err != nil {
			t.Fatalf("generateWithTextImage failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate with pdf file input", func(t *testing.T) {
		buf.Reset()
		err := generateWithPDF(buf)
		if err != nil {
			t.Fatalf("generateWithPDF failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate with video file input (no sound)", func(t *testing.T) {
		buf.Reset()
		err := generateWithMuteVideo(buf)
		if err != nil {
			t.Fatalf("generateWithMuteVideo failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate with video file input", func(t *testing.T) {
		buf.Reset()
		err := generateWithVideo(buf)
		if err != nil {
			t.Fatalf("generateWithVideo failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate with audio file input", func(t *testing.T) {
		buf.Reset()
		err := generateWithAudio(buf)
		if err != nil {
			t.Fatalf("generateWithAudio failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate an audio transcript", func(t *testing.T) {
		buf.Reset()
		err := generateAudioTranscript(buf)
		if err != nil {
			t.Fatalf("generateAudioTranscript failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate with YT video file input", func(t *testing.T) {
		buf.Reset()
		err := generateWithYTVideo(buf)
		if err != nil {
			t.Fatalf("generateWithYTVideo failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate with routing", func(t *testing.T) {
		t.Skip("skipping because of model used in this test. The model 'model-optimizer-exp-04-09' is not consistently available in all test environments.")
		buf.Reset()
		err := generateWithRouting(buf)
		if err != nil {
			t.Fatalf("generateWithRouting failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate chat stream with text prompt", func(t *testing.T) {
		buf.Reset()
		err := generateChatStreamWithText(buf)
		if err != nil {
			t.Fatalf("generateChatStreamWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate Text With PDF", func(t *testing.T) {
		buf.Reset()
		err := generateTextWithPDF(buf)
		if err != nil {
			t.Fatalf("generateTextWithPDF failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate thinking with text prompt", func(t *testing.T) {
		buf.Reset()
		err := generateThinkingWithText(buf)
		if err != nil {
			t.Fatalf("generateThinkingWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate with model optimizer", func(t *testing.T) {
		buf.Reset()
		err := generateModelOptimizerWithTxt(buf)
		if err != nil {
			t.Fatalf("generateModelOptimizerWithTxt failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate chat with text", func(t *testing.T) {
		buf.Reset()
		err := generateChatWithText(buf)
		if err != nil {
			t.Fatalf("generateChatWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	//t.Run("generate text with async stream", func(t *testing.T) {
	//	buf.Reset()
	//	err := generateWithTextAsyncStream(buf)
	//	if err != nil {
	//		t.Fatalf("generateWithTextAsyncStream failed: %v", err)
	//	}
	//
	//	output := buf.String()
	//	if output == "" {
	//		t.Error("expected non-empty output, got empty")
	//	}
	//})

	t.Run("generate with local video file input", func(t *testing.T) {
		buf.Reset()
		err := generateWithLocalVideo(buf)
		if err != nil {
			t.Fatalf("generateWithLocalVideo failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate with local multi local images input", func(t *testing.T) {
		buf.Reset()
		err := generateWithMultiLocalImages(buf)
		if err != nil {
			t.Fatalf("generateWithMultiLocalImages failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

}
