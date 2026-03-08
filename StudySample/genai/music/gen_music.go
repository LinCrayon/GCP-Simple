package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/oto/v2"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// ServerResponse 用来解析 Lyria WebSocket 返回的数据结构
type ServerResponse struct {
	ServerContent *struct {
		AudioChunks []struct {
			Data string `json:"data"` // Base64 编码的音频 PCM 数据
		} `json:"audioChunks"`
	} `json:"serverContent"`
}

// func GenerateMusicStream() {
func GenerateMusicStream() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	// Gemini Lyria 实时音乐生成 WebSocket API
	endpoint := "wss://generativelanguage.googleapis.com/ws/google.ai.generativelanguage.v1alpha.GenerativeService.BidiGenerateMusic"
	url := fmt.Sprintf("%s?key=%s", endpoint, apiKey) //拼接URL

	//拨号,建立WebSocket 连接,请求头放 api key
	conn, _, err := websocket.DefaultDialer.Dial(url, http.Header{
		"x-goog-api-key": []string{apiKey},
		"Content-Type":   []string{"application/json"},
	})
	log.Println("连接成功！")
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	defer conn.Close()

	// 发送 setup 指令 ,告诉服务器使用Lyria音乐模型
	sendJSON(conn, map[string]interface{}{"setup": map[string]interface{}{"model": "models/lyria-realtime-exp"}})

	// 创建 WAV 文件用于保存音频
	outFile, err := os.Create("C:\\Users\\linshengqian\\Desktop\\GCP-Notes\\CodeSample\\GolangSample\\StudySample\\genai\\music\\lyria_music.wav")
	if err != nil {
		log.Println(err)
	}
	defer outFile.Close()

	// WAV编码器：设置录音参数：采样率 (HZ), 每个采样的位深, 声道数,音频编码格式(pcm未压缩 = 1 )
	encoder := wav.NewEncoder(outFile, 48000, 16, 1, 1)
	defer encoder.Close() // WAV 文件必须 Close 才会写 header最终长度,否则无法播放

	// 注册系统信号监听 SIGINT (Ctrl+C) 和 SIGTERM (kill)
	sigChan := make(chan os.Signal, 1)                    //创建channel用来接收操作系统发来的信号，缓冲区大小为 1，没有缓冲区的话信号可能丢失
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM) //注册系统信号监听，当程序收到Interrupt和SIGTERM信号的时候把信号发送到 sigChan

	log.Println("正在生成并录制音频... 按 Ctrl+C 停止并保存")

	// 始化音频播放设备(48000Hz, 单声道, 16位深度 = 2字节)
	samplingRate := 48000
	channelCount := 1
	bitDepthInBytes := 2

	// 创建 oto 音频上下文
	otoCtx, readyChan, err := oto.NewContext(samplingRate, channelCount, bitDepthInBytes)
	if err != nil {
		log.Fatal("初始化音频设备失败:", err)
	}
	// 等待音频设备准备完成
	<-readyChan

	// 创建 Pipe（音频流管道），Pipe 是一个内存管道， pw.Write(data)  --->  pr.Read(data)
	// Google 音频数据 → pw → pr → oto 播放器
	pr, pw := io.Pipe()

	// 创建播放器
	player := otoCtx.NewPlayer(pr)
	defer player.Close()

	// 启动播放器, Player会不断从 pr 读取音频并播放
	go func() {
		player.Play()
	}()

	// 启动 Goroutine 接收 WebSocket 音频流
	go func() {
		for {
			//读取 WebSocket 消息
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
			}
			//fmt.Println(string(message))
			//判断模型是否初始化完成
			if strings.Contains(string(message), "setupComplete") {
				log.Println("模型初始化完成")
				// 指令 B: 设置音乐风格
				sendJSON(conn, map[string]interface{}{
					"clientContent": map[string]interface{}{
						"weightedPrompts": []interface{}{
							map[string]interface{}{
								"text":   "minimal techno",
								"weight": 1.0,
							},
						},
					},
				})
				// 指令 C: 设置 BPM,设置节奏为 90 BPM (每分钟 90 拍)
				sendJSON(conn, map[string]interface{}{"musicGenerationConfig": map[string]interface{}{"bpm": 90, "temperature": 1.0}})
				// 指令 D: 开始播放
				sendJSON(conn, map[string]interface{}{
					"playbackControl": "PLAY",
				})
				continue
			}

			// 解析服务器返回 JSON
			var resp ServerResponse
			if err := json.Unmarshal(message, &resp); err == nil {
				if resp.ServerContent != nil {
					for _, chunk := range resp.ServerContent.AudioChunks {
						// Base64 → 原始 PCM byte
						rawBytes, _ := base64.StdEncoding.DecodeString(chunk.Data)

						// 直接将字节写入播放器Pipe，它会自动传给声卡
						// 将解码后的字节写入管道, Player 会自动从 pr 中读取这些字节并传给声卡
						// pw.Write -> pr.Read -> oto 播放
						_, err := pw.Write(rawBytes)
						if err != nil {
							log.Println("写入播放管道失败:", err)
						}

						// 将 byte 转换为 wav 库需要的 int 采样数据
						// 因为是 16-bit (2字节)，我们将每两个字节合成为一个 int16
						// WAV 库不认识原始字节，它只认识数字,调用 bytesToInts 把每 2 个字节变成 1 个数字
						buf := &audio.IntBuffer{
							Data:           bytesToInts(rawBytes),
							Format:         &audio.Format{SampleRate: 48000, NumChannels: 1},
							SourceBitDepth: 16,
						}
						// 写进录音带
						if err := encoder.Write(buf); err != nil {
							log.Println("写入错误:", err)
						}
						fmt.Printf("🎵 音频块: %d bytes\n", len(rawBytes))
					}
				}
			}
		}
	}()

	//主线程阻塞等待 Ctrl+C
	<-sigChan
	fmt.Println("\n正在结束录制并更新文件头...")
}

// 把 Go 对象编码为 JSON 并通过 WebSocket 发送
func sendJSON(c *websocket.Conn, v interface{}) {
	data, _ := json.Marshal(v)
	c.WriteMessage(websocket.TextMessage, data)
}

// 把原始字节 (Byte) 转换成数字 (Int) ,就像把一堆乐谱符号转换成钢琴按键的力度数值
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

/*
 //TODO 架构
Gemini Lyria
      │
      │ WebSocket
      ▼
Base64 音频
      │
      ▼
Decode PCM
      │
      ▼
      ┌───────────────┐
      │               │
      ▼               ▼
 Pipe 写入        WAV Encoder
      │
      ▼
 oto Player
      │
      ▼
 声卡播放
*/

/*
TODO 执行流程
启动程序
   │
连接 WebSocket
   │
发送 setup
   │
等待 setupComplete
   │
发送 prompt + bpm
   │
发送 PLAY
   │
服务器返回音频 chunk
   │
Base64 decode
   │
写入 pipe 播放
   │
同时写入 WAV
   │
Ctrl+C
   │
encoder.Close()
   │
WAV 文件完成
*/
