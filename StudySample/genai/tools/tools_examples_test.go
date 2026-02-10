package tools

import (
	"bytes"
	"testing"
)

func TestTextGeneration(t *testing.T) {
	//tc := testutil.SystemTest(t)
	//
	//t.Setenv("GOOGLE_GENAI_USE_VERTEXAI", "1")
	//t.Setenv("GOOGLE_CLOUD_LOCATION", "us-central1")
	//t.Setenv("GOOGLE_CLOUD_PROJECT", tc.ProjectID)

	buf := new(bytes.Buffer)

	t.Run("使用函数声明和函数响应生成", func(t *testing.T) {
		buf.Reset()
		err := generateWithFuncCall(buf)
		if err != nil {
			t.Fatalf("generateWithFuncCall failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

}
