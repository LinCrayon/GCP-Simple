package music

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/gorilla/websocket"
)

// ServerResponse 定义 API 返回结构 定义了 Google 传回来的消息长什么样
type ServerResponse struct {
	ServerContent struct {
		AudioChunks []struct {
			Data string `json:"data"` // 这里的 Data 就是 Base64 编码的音频片段
		} `json:"audioChunks"`
	} `json:"serverContent"`
}

func GenerateMusicStream() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	endpoint := "wss://generativelanguage.googleapis.com/ws/google.ai.generativelanguage.v1alpha.GenerativeService.BidiGenerateMusic"
	url := fmt.Sprintf("%s?key=%s", endpoint, apiKey)

	// 1. 拨号连接,我们在请求头（Header）里放上 API Key，证明我们是合法用户
	conn, _, err := websocket.DefaultDialer.Dial(url, http.Header{
		"x-goog-api-key": []string{apiKey},
		"Content-Type":   []string{"application/json"},
	})
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	defer conn.Close()

	// 发送指令
	// 指令 A: 用 Lyria 这个音乐模型
	sendJSON(conn, map[string]interface{}{"setup": map[string]interface{}{"model": "models/lyria-realtime-exp"}})
	// 指令 B: 听 minimal techno 风格的音乐
	sendJSON(conn, map[string]interface{}{"clientContent": map[string]interface{}{"weightedPrompts": []interface{}{map[string]interface{}{"text": "minimal techno", "weight": 1.0}}}})
	// 指令 C: 设置节奏为 90 BPM (每分钟 90 拍)
	sendJSON(conn, map[string]interface{}{"musicGenerationConfig": map[string]interface{}{"bpm": 90, "temperature": 1.0}})
	// 指令 D: 开始播放！
	sendJSON(conn, map[string]interface{}{"playbackControl": "PLAY"})

	// 2. 创建 WAV 文件准备写入
	outFile, err := os.Create("C:\\Users\\linshengqian\\Desktop\\GCP-Notes\\CodeSample\\GolangSample\\StudySample\\genai\\music\\lyria_music.wav")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	// 设置录音参数：采样率 48000Hz (高保真), 16位深度, 1个声道 (单声道)
	encoder := wav.NewEncoder(outFile, 48000, 16, 1, 1)
	defer encoder.Close()

	// 3. 监听“停止”信号 (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	fmt.Println("正在生成并录制音频... 按 Ctrl+C 停止并保存")

	// 4. 开启一个“后台窗口” (Goroutine) 来不断接收音频块
	go func() {
		for {
			// 读取服务器发来的消息
			_, message, err := conn.ReadMessage()
			if err != nil {
				return
			}

			var resp ServerResponse
			if err := json.Unmarshal(message, &resp); err == nil {
				for _, chunk := range resp.ServerContent.AudioChunks {
					//base64 ----> byte
					rawBytes, _ := base64.StdEncoding.DecodeString(chunk.Data)

					// 将 byte 转换为 wav 库需要的 int 采样数据
					// 因为是 16-bit (2字节)，我们将每两个字节合成为一个 int16
					// WAV 库不认识原始字节，它只认识数字
					// 我们调用 bytesToInts 把每 2 个字节变成 1 个数字
					buf := &audio.IntBuffer{
						Data:           bytesToInts(rawBytes),
						Format:         &audio.Format{SampleRate: 48000, NumChannels: 1},
						SourceBitDepth: 16,
					}
					// 把数字写进录音带
					if err := encoder.Write(buf); err != nil {
						log.Println("写入错误:", err)
					}
					fmt.Printf("已录制音频块: %d 字节\n", len(rawBytes))
				}
			}
		}
	}()

	//阻塞在这里，直到你按 Ctrl+C
	<-sigChan
	fmt.Println("\n正在结束录制并更新文件头...")
}

// 辅助工具：把 Go 的对象转成 JSON 字符串并发给 Google
func sendJSON(c *websocket.Conn, v interface{}) {
	data, _ := json.Marshal(v)
	c.WriteMessage(websocket.TextMessage, data)
}

// 辅助工具：把原始字节 (Byte) 转换成数字 (Int)
// 就像把一堆乐谱符号转换成钢琴按键的力度数值
func bytesToInts(b []byte) []int {
	// 16位音频意味着 2 个字节代表 1 个声音采样点，所以长度减半
	n := len(b) / 2
	ints := make([]int, n)
	for i := 0; i < n; i++ {
		// 这里采用 Little-endian (小端序) 算法
		// 将两个 byte 拼成一个 int16 数字
		ints[i] = int(int16(b[i*2]) | int16(b[i*2+1])<<8)
	}
	return ints
}
