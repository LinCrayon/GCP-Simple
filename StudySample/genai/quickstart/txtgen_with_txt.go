package quickstart

import (
	"context"
	"fmt"
	"google.golang.org/genai"
	"io"
)

/**
 * @Description 提出你的第一个请求
 * @Author linshengqian
 * @Date 2026/2/10 15:34
 **/
func generateWithText(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("创建genai客户端失败: %w", err)
	}

	resp, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text("人工智能是如何运作的？"),
		nil)
	if err != nil {
		return fmt.Errorf("生成内容失败: %w", err)
	}

	respText := resp.Text()
	fmt.Fprintln(w, respText)

	return nil
}

//func main() {
//	buf := new(bytes.Buffer)
//	err := generateWithText(buf)
//	if err != nil {
//		log.Fatalf("generateWithText failed: %v", err)
//	}
//	output := buf.String()
//	if output == "" {
//		log.Fatal("expected non-empty output")
//	}
//	log.Printf("model output: %s", output)
//}
