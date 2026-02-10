package quickstart

import (
	"bytes"
	"os"
	"testing"
)

func TestGenerate(t *testing.T) {
	//t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	buf := new(bytes.Buffer)

	t.Run("提出你的第一个请求", func(t *testing.T) {
		buf.Reset()
		err := generateWithText(buf)
		if err != nil {
			t.Fatalf("generateWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
		//t.Log(output) //测试的时候加 go test -v
	})

	t.Run("生成图片", func(t *testing.T) {
		err := generateMMFlashWithText(buf)
		if err != nil {
			t.Fatalf("generateMMFlashWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Fatal("expected some text output, got empty")
		}

		// 检查是否生成了图片文件
		const imageFile = "example-image-eiffel-tower.png"
		if _, err := os.Stat(imageFile); err != nil {
			t.Fatalf("expected image file %s to exist: %v", imageFile, err)
		}
	})

	t.Run("图像识别-识别 1.png 图片内容", func(t *testing.T) {
		buf.Reset()

		err := generateWithTextImage(buf)
		if err != nil {
			t.Fatalf("生成图像识别失败: %v", err)
		}

		output := buf.String()

		// 检查是否返回了文本内容
		if output == "" {
			t.Fatal("expected some text output, got empty")
		}

		t.Logf("模型输出: %s", output)
	})

	t.Run("代码执行-计算斐波那契并找回文", func(t *testing.T) {
		buf.Reset()

		err := generateWithCodeExec(buf)
		if err != nil {
			t.Fatalf("GenerateWithCodeExec failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Fatal("expected some text output, got empty")
		}

		t.Logf("模型输出:\n%s", output)

		// 可选：检查是否包含关键字，例如 Fibonacci 或 Outcome_OK
		if !containsKeyword(output, "fib_20") {
			t.Errorf("output does not contain expected Fibonacci result: %s", output)
		}
		if !containsKeyword(output, "Outcome: OUTCOME_OK") {
			t.Errorf("code execution did not succeed: %s", output)
		}
	})

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

// 辅助函数：检查文本中是否包含关键字
func containsKeyword(text, keyword string) bool {
	return bytes.Contains([]byte(text), []byte(keyword))
}
