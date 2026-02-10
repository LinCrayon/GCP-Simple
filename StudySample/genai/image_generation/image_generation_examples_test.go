package image_generation

import (
	"bytes"
	"os"
	"testing"
)

// 集成测试：依赖外部资源
func TestGenerateMMFlashWithText(t *testing.T) {
	// 跳过短测试（避免 CI 默认跑）
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	var buf bytes.Buffer

	err := generateMMFlashWithText(&buf)
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

	// 清理文件（避免污染工作区）
	//if err := os.Remove(imageFile); err != nil {
	//	t.Logf("warning: failed to remove image file: %v", err)
	//}
}

// TODO 单元测试：偏向快速验证逻辑（只关心函数输出，不依赖外部资源）图片可能生成了，但系统不关心
func TestImageGeneration(t *testing.T) {
	//tc := testutil.SystemTest(t)
	//
	//t.Setenv("GOOGLE_GENAI_USE_VERTEXAI", "1")
	//t.Setenv("GOOGLE_CLOUD_LOCATION", "us-central1")
	//t.Setenv("GOOGLE_CLOUD_PROJECT", tc.ProjectID)

	buf := new(bytes.Buffer)

	t.Run("generate multimodal flash content with text and image", func(t *testing.T) {
		buf.Reset()
		err := generateMMFlashWithText(buf)
		if err != nil {
			t.Fatalf("generateMMFlashWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("生成包含文本和图像的多模式 Flash 内容", func(t *testing.T) {
		buf.Reset()
		err := generateMMFlashTxtImgWithText(buf)
		if err != nil {
			t.Fatalf("generateMMFlashTxtImgWithText failed: %v", err)
		}

		// 函数预期会向 writer 输出生成的 Markdown 文件路径
		output := buf.String()
		if output == "" {
			t.Error("期望有输出内容，但得到的是空字符串")
		}
	})

}
