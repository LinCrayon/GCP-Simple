package live

// [START googlegenaisdk_live_txt_with_audio]
import (
	"context"
	"fmt"
	"google.golang.org/genai"
	"io"
	"net/http"
	"os"
)

// generateLiveTextWithAudio 音频输入，文本输出
// 发送音频，接收文本回复
func generateLiveTextWithAudio(ctx context.Context) error {
	///////////////////////////////////////
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	//os.Setenv("HTTPS_PROXY", "socks5://127.0.0.1:10606")
	//os.Setenv("HTTP_PROXY", "socks5://127.0.0.1:10606")
	//
	//// 确保程序退出时清理（可选）
	//defer os.Unsetenv("HTTPS_PROXY")

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Location: "us-central1",
		Project:  "train-crayon-20260304",
		Backend:  genai.BackendVertexAI,
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	// 4. 后续直接使用 SDK 的 Connect 方法即可
	modelName := "gemini-live-2.5-flash-preview-native-audio-09-2025"
	config := &genai.LiveConnectConfig{
		ResponseModalities: []genai.Modality{genai.ModalityAudio},
	}

	session, err := client.Live.Connect(ctx, modelName, config)
	if err != nil {
		// 如果这里报错，SDK 会自动使用上面配置的 httpClient (及其代理)
		return fmt.Errorf("failed to connect live: %w", err)
	}
	defer session.Close()

	audioURL := "https://storage.googleapis.com/generativeai-downloads/data/16000.wav"
	// Download audio
	resp, err := http.Get(audioURL)
	if err != nil {
		return fmt.Errorf("failed to download audio: %w", err)
	}
	defer resp.Body.Close()

	audioBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read audio: %w", err)
	}

	//fmt.Fprintf(w, "> Answer to this audio url: %s\n\n", audioURL)

	// Send the audio as Blob media input
	err = session.SendRealtimeInput(genai.LiveRealtimeInput{
		Media: &genai.Blob{
			Data:     audioBytes,
			MIMEType: "audio/pcm;rate=16000",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send audio input: %w", err)
	}

	// Stream the response
	//var response strings.Builder
	var audioData []byte

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Timeout reached, ending session")
			goto END
		default:
			chunk, err := session.Receive()
			if err != nil {
				if err == io.EOF {
					break
				}
				return fmt.Errorf("error receiving response: %w", err)
			}

			if chunk.ServerContent == nil {
				continue
			}

			// Handle model turn responses
			if chunk.ServerContent.ModelTurn != nil {
				for _, part := range chunk.ServerContent.ModelTurn.Parts {
					if part.InlineData != nil {
						audioData = append(audioData, part.InlineData.Data...)
						fmt.Printf("received audio chunk: %d bytes\n", len(part.InlineData.Data))
					}
				}
			}
		}
	}
END:
	fmt.Println("Ending session, total audio bytes:", len(audioData))
	err = os.WriteFile("gemini_output.raw", audioData, 0644)
	if err != nil {
		return err
	}

	err = writeWav("gemini_output.wav", audioData, 24000)
	if err != nil {
		return err
	}

	fmt.Println("WAV file saved: gemini_output.wav")

	// Example output:
	// > Answer to this audio url: https://storage.googleapis.com/generativeai-downloads/data/16000.wav
	// Yes, I can hear you. How can I help you today?
	return nil
}

func writeWav(filename string, pcmData []byte, sampleRate int) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	dataSize := uint32(len(pcmData))
	byteRate := uint32(sampleRate * 2) // 16bit mono

	header := []byte{
		'R', 'I', 'F', 'F',
		byte(dataSize + 36), byte((dataSize + 36) >> 8), byte((dataSize + 36) >> 16), byte((dataSize + 36) >> 24),
		'W', 'A', 'V', 'E',
		'f', 'm', 't', ' ',
		16, 0, 0, 0,
		1, 0,
		1, 0,
		byte(sampleRate), byte(sampleRate >> 8), byte(sampleRate >> 16), byte(sampleRate >> 24),
		byte(byteRate), byte(byteRate >> 8), byte(byteRate >> 16), byte(byteRate >> 24),
		2, 0,
		16, 0,
		'd', 'a', 't', 'a',
		byte(dataSize), byte(dataSize >> 8), byte(dataSize >> 16), byte(dataSize >> 24),
	}

	_, err = f.Write(header)
	if err != nil {
		return err
	}

	_, err = f.Write(pcmData)
	return err
}

//func init() {
//	proxyAddr := "127.0.0.1:10606"
//	dialer, _ := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
//
//	// 修改全局默认传输层
//	http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
//		return dialer.Dial(network, addr)
//	}
//}
