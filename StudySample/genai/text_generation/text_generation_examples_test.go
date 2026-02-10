package quickstart

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

}
