package live

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateLiveTranscribeWithAudio 音频转录：处理音频输入，进行实时语音识别和文本回复
func generateLiveTranscribeWithAudio(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("创建 genai 客户端失败: %w", err)
	}

	modelName := "gemini-live-2.5-flash-preview-native-audio"

	// 实时连接配置
	config := &genai.LiveConnectConfig{
		ResponseModalities:       []genai.Modality{genai.ModalityAudio}, // 输出和输入都支持音频
		InputAudioTranscription:  &genai.AudioTranscriptionConfig{},     // 配置输入音频转录
		OutputAudioTranscription: &genai.AudioTranscriptionConfig{},     // 配置输出音频转录
	}

	// 打开实时会话
	session, err := client.Live.Connect(ctx, modelName, config)
	if err != nil {
		return fmt.Errorf("failed to connect live session: %w", err)
	}
	defer session.Close()

	// 模拟用户输入（文本或可对应音频内容）
	inputText := "你好？Gemini，你在吗？"
	fmt.Fprintf(w, "> %s\n", inputText)

	// 发送客户端内容
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
		return fmt.Errorf("发送客户端内容失败: %w", err)
	}

	// 拼接模型输出文本
	var response string

	// 实时接收模型返回
	for {
		message, err := session.Receive()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("接收实时流出错: %w", err)
		}

		if message.ServerContent != nil {
			// 模型生成的完整 Turn
			if message.ServerContent.ModelTurn != nil {
				fmt.Fprintf(w, "Model turn: %v\n", message.ServerContent.ModelTurn)
			}

			// 用户音频输入转录结果
			if message.ServerContent.InputTranscription != nil {
				if message.ServerContent.InputTranscription.Text != "" {
					fmt.Fprintf(w, "用户语音转录文本: %s\n",
						message.ServerContent.InputTranscription.Text)
				}
			}

			// 模型输出转录文本
			if message.ServerContent.OutputTranscription != nil {
				if message.ServerContent.OutputTranscription.Text != "" {
					response += message.ServerContent.OutputTranscription.Text
				}
			}
		}
	}
	// 示例输出：
	// > 你好？Gemini，你在吗？
	// 我在。你想聊些什么呢？
	//实时音频交互
	//用户音频自动转文本（InputTranscription）
	//模型音频生成可转文本（OutputTranscription）
	fmt.Fprintln(w, response)
	return nil
}
