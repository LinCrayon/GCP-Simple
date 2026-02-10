package quickstart

import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

/**
 * @Description
 * @Author linshengqian
 * @Date 2026/2/10 16:40
 **/
// generateWithCodeExec 展示如何使用代码执行工具生成文本
func generateWithCodeExec(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	//用户输入
	prompt := "计算第 20 个斐波那契数。然后找到离它最近的回文"
	// GenAI 使用对话式输入，每个 Content 对象可以包含多个 Part
	contents := []*genai.Content{
		{Parts: []*genai.Part{
			{Text: prompt}, // 用户的问题或指令
		},
			Role: genai.RoleUser},
	}
	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{ //指定工具
			{CodeExecution: &genai.ToolCodeExecution{}}, // 启用代码执行工具
		},
		// 温度设置 0.0 表示尽量确定性的回答，不引入随机性
		Temperature: genai.Ptr(float32(0.0)),
	}
	modelName := "gemini-2.5-flash"

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	//遍历返回内容的 Parts
	for _, p := range resp.Candidates[0].Content.Parts {
		if p.Text != "" { //输出文本
			fmt.Fprintf(w, "Gemini: %s", p.Text)
		}
		if p.ExecutableCode != nil { //输出模型生成的可执行代码
			fmt.Fprintf(w, "Language: %s\n%s\n", p.ExecutableCode.Language, p.ExecutableCode.Code)
		}
		if p.CodeExecutionResult != nil { //输出代码执行结果
			fmt.Fprintf(w, "Outcome: %s\n%s\n", p.CodeExecutionResult.Outcome, p.CodeExecutionResult.Output)
		}
	}

	// Example response:
	// Gemini: Okay, I can do that. First, I'll calculate the 20th Fibonacci number. Then, I need ...
	//
	// Language: PYTHON
	//
	// def fibonacci(n):
	//    ...
	//
	// fib_20 = fibonacci(20)
	// print(f'{fib_20=}')
	//
	// Outcome: OUTCOME_OK
	// fib_20=6765
	//
	// Now that I have the 20th Fibonacci number (6765), I need to find the nearest palindrome. ...
	// ...

	return nil
}
