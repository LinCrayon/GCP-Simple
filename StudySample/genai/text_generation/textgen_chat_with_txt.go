package text_generation

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateChatWithText 多轮对话：保持上下文的对话、记忆历史消息
func generateChatWithText(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("创建 genai 客户端失败： %w", err)
	}
	modelName := "gemini-2.5-flash"
	history := []*genai.Content{
		{
			Role: genai.RoleUser,
			Parts: []*genai.Part{
				{Text: "你好呀"},
			},
		},
		{
			Role: genai.RoleModel,
			Parts: []*genai.Part{
				{Text: "很高兴认识你。您想知道什么？"},
			},
		},
	}
	chatSession, err := client.Chats.Create(ctx, modelName, nil, history)
	if err != nil {
		return fmt.Errorf("无法创建 genai 聊天会话： %w", err)
	}
	contents := genai.Part{Text: "给我讲个故事吧。"}
	resp, err := chatSession.SendMessage(ctx, contents)
	if err != nil {
		return fmt.Errorf("发送消息失败： %w", err)
	}

	respText := resp.Text()

	fmt.Fprintln(w, respText)
	// 响应示例：
	// 好吧，安顿下来。让我告诉你一个关于一位安静的制图师的故事，但不是关于陆地和海洋的故事。
	// ...
	// 在幽静的奥克黑文镇，坐落在低语山和低语河之间，住着一位名叫埃拉拉的女人。
	// ...

	return nil
}
