package text_generation

import (
	"context"
	"fmt"
	"google.golang.org/genai"
	"io"
)

/**
 * @Description
 * @Author linshengqian
 * @Date 2026/2/10 17:58
 **/
func generateWithTextStream(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.5-flash"
	contents := genai.Text("天空为什么是蓝色的？")

	//流驱动函数
	stream := client.Models.GenerateContentStream(ctx, modelName, contents, nil)

	var streamErr error

	//stream 接收一个回调函数（yield）: true  -> 继续接收下一个 chunk, false -> 立即终止流式生成
	stream(func(resp *genai.GenerateContentResponse, err error) bool {
		if err != nil {
			streamErr = err
			return false
		}
		fmt.Printf("收到 chunk: %q\n", resp.Text()) //测试查看流块
		fmt.Fprintln(w, resp.Text())              //当前 块 的文本内容
		return true
	})

	if streamErr != nil {
		return fmt.Errorf("failed to generate content: %w", streamErr)
	}

	return nil
}
