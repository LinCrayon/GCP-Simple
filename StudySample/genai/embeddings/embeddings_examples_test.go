package embeddings

import (
	"bytes"
	"testing"
)

func TestEmbedGeneration(t *testing.T) {
	////tc := testutil.SystemTest(t)
	//
	//t.Setenv("GOOGLE_GENAI_USE_VERTEXAI", "1")
	//t.Setenv("GOOGLE_CLOUD_LOCATION", "us-central1")
	//t.Setenv("GOOGLE_CLOUD_PROJECT", tc.ProjectID)

	buf := new(bytes.Buffer)

	t.Run("生成带有文本的嵌入内容", func(t *testing.T) {
		buf.Reset()
		err := generateEmbedContentWithText(buf)
		if err != nil {
			t.Fatalf("generateEmbedContentWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

}
