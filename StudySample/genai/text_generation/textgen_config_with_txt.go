package text_generation

import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateWithConfig 演示如何使用【文本提示 + 自定义配置】生成文本内容。
// 本示例展示了如何设置温度、候选数量以及返回的 MIME 类型。
func generateWithConfig(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("创建 genai 客户端失败: %w", err)
	}

	modelName := "gemini-2.5-flash"

	contents := genai.Text("为什么天空是蓝色的？")

	// 生成配置： https://pkg.go.dev/google.golang.org/genai#GenerateContentConfig
	config := &genai.GenerateContentConfig{
		Temperature:      genai.Ptr(float32(0.0)), // 温度设为 0，保证输出稳定、确定
		CandidateCount:   int32(1),                // 只返回一个候选结果
		ResponseMIMEType: "application/json",      // 指定返回内容的 MIME 类型为 JSON
	}

	// 调用模型生成内容
	resp, err := client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		return fmt.Errorf("生成内容失败: %w", err)
	}

	// 提取模型返回的文本结果
	respText := resp.Text()

	// 将生成结果输出到 writer
	fmt.Fprintln(w, respText)

	// 示例输出：
	// {
	//   "explanation": "天空呈现蓝色是由于一种称为瑞利散射的现象……"
	// }

	return nil
}
