package tools

import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateWithCodeExec 演示如何使用「代码执行（Code Execution）」工具生成文本。
// 模型会自动生成可执行代码、运行代码，并基于执行结果给出最终回答。
func generateWithCodeExec(w io.Writer) error {

	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("创建 genai 客户端失败: %w", err)
	}

	prompt := "计算第 20 个斐波那契数，然后找出最接近它的回文数。"

	contents := []*genai.Content{
		{
			Role: genai.RoleUser,
			Parts: []*genai.Part{
				{Text: prompt},
			},
		},
	}

	// 生成配置： 启用 Code Execution 工具
	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{
				CodeExecution: &genai.ToolCodeExecution{},
			},
		},
		// 温度设为 0，保证推理与代码生成稳定
		Temperature: genai.Ptr(float32(0.0)),
	}

	modelName := "gemini-2.5-flash"

	// 调用模型生成内容
	resp, err := client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		return fmt.Errorf("生成内容失败: %w", err)
	}

	// 遍历模型返回的内容片段
	for _, p := range resp.Candidates[0].Content.Parts {
		// 模型生成的自然语言说明
		if p.Text != "" {
			fmt.Fprintf(w, "Gemini 输出：%s", p.Text)
		}

		// 模型生成的可执行代码
		if p.ExecutableCode != nil {
			fmt.Fprintf(
				w,
				"代码语言：%s\n%s\n",
				p.ExecutableCode.Language,
				p.ExecutableCode.Code,
			)
		}

		// 代码执行结果
		if p.CodeExecutionResult != nil {
			fmt.Fprintf(
				w,
				"执行结果状态：%s\n执行输出：\n%s\n",
				p.CodeExecutionResult.Outcome,
				p.CodeExecutionResult.Output,
			)
		}
	}

	// 示例输出：
	// Gemini 输出：好的，我可以先计算第 20 个斐波那契数，然后寻找最接近的回文数……
	//
	// 代码语言：PYTHON
	//
	// def fibonacci(n):
	//     ...
	//
	// fib_20 = fibonacci(20)
	// print(f'{fib_20=}')
	//
	// 执行结果状态：OUTCOME_OK
	// 执行输出：
	// fib_20=6765
	//
	// 接下来我会基于 6765 寻找最接近的回文数……
	return nil
}
