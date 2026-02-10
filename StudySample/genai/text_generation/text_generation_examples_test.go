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

}
