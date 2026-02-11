package live

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateLiveWithText 演示如何使用 Gemini 实时模型（Live Model）
// 进行文本交互，并实时接收模型回复。
func generateLiveWithText(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("创建 genai 客户端失败: %w", err)
	}

	modelName := "gemini-2.0-flash-live-preview-04-09"

	// 实时连接配置：只使用文本模式
	config := &genai.LiveConnectConfig{
		ResponseModalities: []genai.Modality{genai.ModalityText},
	}

	// 打开实时会话
	session, err := client.Live.Connect(ctx, modelName, config)
	if err != nil {
		return fmt.Errorf("实时连接失败: %w", err)
	}
	defer session.Close()

	// 准备用户输入
	inputText := "你好？Gemini，你在吗？"
	fmt.Fprintf(w, "> %s\n\n", inputText)

	// Send text content to the model
	err = session.SendClientContent(genai.LiveClientContentInput{
		Turns: []*genai.Content{
			{
				Role: genai.RoleUser,
				Parts: []*genai.Part{
					{Text: inputText},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("发送内容失败: %w", err)
	}

	// 实时接收模型回复
	var response string
	for {
		chunk, err := session.Receive()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("接收回复出错: %w", err)
		}

		if chunk.ServerContent != nil && chunk.ServerContent.ModelTurn != nil {
			for _, part := range chunk.ServerContent.ModelTurn.Parts {
				if part.Text != "" {
					response += part.Text
				}
			}
		}
	}

	// 示例输出：
	// > 你好？Gemini，你在吗？
	// 我在。你想聊些什么呢？
	fmt.Fprintln(w, response)
	return nil
}
