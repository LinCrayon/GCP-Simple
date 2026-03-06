package live

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()

	// 创建 Vertex AI client
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:  "train-crayon-20260304",
		Location: "us-central1",
		Backend:  genai.BackendVertexAI,
	})
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	model := "gemini-live-2.5-flash-preview-native-audio-09-2025"

	config := &genai.LiveConnectConfig{
		ResponseModalities: []genai.Modality{
			genai.ModalityAudio,
		},
	}

	// 建立 Live WebSocket 会话
	session, err := client.Live.Connect(ctx, model, config)
	if err != nil {
		log.Fatalf("failed to connect live session: %v", err)
	}
	defer session.Close()

	fmt.Println("Session established. Ready to send audio...")

	// 这里只是示例建立连接
	// 实际项目通常会 SendRealtimeInput()

	select {}
}
