package main

import (
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/genai"
	"os"
	"strings"
	"time"
)

// TranscriptionResult 对应 Gemini 返回的最外层 JSON
type TranscriptionResult struct {
	Summary  string    `json:"summary"`  // 视频 / 音频整体摘要
	Segments []Segment `json:"segments"` // 分段转录结果
}

// Segment 表示一段音频的转录信息
type Segment struct {
	Timestamp    []string `json:"timestamp"`             // 时间戳，如 ["00:01", "00:05"]
	Content      string   `json:"content"`               // 原始识别文本
	Language     string   `json:"language"`              // 语言名称（如 English）
	LanguageCode string   `json:"language_code"`         // 语言代码（如 en / zh-CN）
	Translation  string   `json:"translation,omitempty"` // 翻译文本（如果不是中文）
	Emotion      string   `json:"emotion"`               // 情绪（Gemini 推断）
}

// VideoAudioToJSON 1：视频 / 音频 → 转录 JSON（Gemini）
//   - 输入 YouTube / 音频 URL
//   - 使用 Gemini 多模态能力
//   - 输出结构化 JSON 转录结果
//   - 如果 videoURL 是 gs://，则必须改用 Vertex AI client
func VideoAudioToJSON(ctx context.Context, client *genai.Client, videoURL string) (string, error) {
	prompt := `
处理音频并生成详细的转录。
将非中文片段翻译成中文。
Return JSON with:
- summary
- segments (timestamp, content, language_code, translation, emotion)
`
	var resp *genai.GenerateContentResponse
	// 使用指数退避的重试机制，避免 503 / UNAVAILABLE
	err := retryWithBackoff(ctx, 5, func() error {
		var err error
		resp, err = client.Models.GenerateContent(
			ctx,
			"gemini-3-flash-preview",
			[]*genai.Content{
				{
					Parts: []*genai.Part{
						{
							// FileData：传入视频或音频 URL
							// Gemini API 仅支持 http(s)，不支持 gs://
							FileData: &genai.FileData{
								FileURI: videoURL,
							},
						},
						{
							Text: prompt,
						},
					},
				},
			},
			nil,
		)
		return err
	})
	if err != nil {
		return "", err
	}

	return resp.Text(), nil
}

// TranslateToChinese JSON → 中文文本
// 功能：
//   - 清洗 Gemini 返回的 JSON
//   - 反序列化为结构体
//   - 抽取中文内容（或翻译内容）
//   - 合并成一段可朗读文本
func TranslateToChinese(jsonStr string) (string, error) {
	// Gemini 有时会包 ```json ```，需要清洗
	jsonStr = cleanJSON(jsonStr)

	var result TranscriptionResult

	// JSON → Go struct
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return "", err
	}

	if len(result.Segments) == 0 {
		return "", fmt.Errorf("转录结果为空")
	}

	var chineseTexts []string

	for _, seg := range result.Segments {
		// 已经是中文
		if strings.HasPrefix(seg.LanguageCode, "zh") {
			chineseTexts = append(chineseTexts, seg.Content)
			continue
		}

		// Gemini 已给翻译
		if seg.Translation != "" {
			chineseTexts = append(chineseTexts, seg.Translation)
			continue
		}

		// 兜底
		chineseTexts = append(chineseTexts, seg.Content)
	}
	// 用中文句号拼接，适合 TTS
	return strings.Join(chineseTexts, "。"), nil
}

// ChineseTextToSpeechLong 长文本 TTS（自动切片）
// 功能：
//   - 将超长中文文本切片
//   - 多次调用 TTS
//   - 拼接成一个 MP3 文件
//
// 原因：
//   - Cloud TTS 对单次文本长度有限制
func ChineseTextToSpeechLong(ctx context.Context, text string, outputFile string) error {
	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	//切片文本
	chunks := splitTextByBytes(text, 4800)

	var finalAudio []byte

	for i, chunk := range chunks {
		fmt.Printf("合成第 %d / %d 段\n", i+1, len(chunks))

		req := &texttospeechpb.SynthesizeSpeechRequest{
			//要朗读的内容
			Input: &texttospeechpb.SynthesisInput{
				InputSource: &texttospeechpb.SynthesisInput_Text{
					Text: chunk,
				},
			},
			//朗读的声音（语言、音色）
			Voice: &texttospeechpb.VoiceSelectionParams{
				LanguageCode: "cmn-CN",           //普通话大陆
				Name:         "cmn-CN-Wavenet-A", //A女声, Wavenet：高质量神经网络语音
			},
			//输出音频的格式和参数
			AudioConfig: &texttospeechpb.AudioConfig{
				AudioEncoding: texttospeechpb.AudioEncoding_MP3, //输出 MP3 格式
			},
		}

		resp, err := client.SynthesizeSpeech(ctx, req)
		if err != nil {
			return fmt.Errorf("TTS 第 %d 段失败: %w", i+1, err)
		}

		finalAudio = append(finalAudio, resp.AudioContent...)
	}

	return os.WriteFile(outputFile, finalAudio, 0644)
}

// 数据清洗
func cleanJSON(raw string) string {
	s := strings.TrimSpace(raw)

	// 去掉 ```json
	if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```")
		s = strings.TrimPrefix(s, "json")
		s = strings.TrimSpace(s)
	}

	// 去掉结尾 ```
	if strings.HasSuffix(s, "```") {
		s = strings.TrimSuffix(s, "```")
		s = strings.TrimSpace(s)
	}

	// 只保留最外层 JSON
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start >= 0 && end > start {
		s = s[start : end+1]
	}

	return s
}

// 按 byte 切中文文本 : 按字节数切分字符串（避免 TTS 超限），每一段都不超过 maxBytes，用于避免 TTS 文本超长报错。
func splitTextByBytes(text string, maxBytes int) []string {
	var chunks []string
	var current strings.Builder
	currentBytes := 0

	for _, r := range text {
		b := len([]byte(string(r)))
		if currentBytes+b > maxBytes {
			chunks = append(chunks, current.String())
			current.Reset()
			currentBytes = 0
		}
		current.WriteRune(r)
		currentBytes += b
	}

	if current.Len() > 0 {
		chunks = append(chunks, current.String())
	}

	return chunks
}

// 通用重试函数
// - 对 503 / UNAVAILABLE 错误进行指数退避重试
func retryWithBackoff(
	ctx context.Context,
	maxRetry int,
	fn func() error,
) error {

	delay := time.Second

	for i := 0; i < maxRetry; i++ {
		err := fn()
		if err == nil {
			return nil
		}

		// 只对 503 / UNAVAILABLE 重试
		if !strings.Contains(err.Error(), "UNAVAILABLE") &&
			!strings.Contains(err.Error(), "503") {
			return err
		}

		fmt.Printf("服务繁忙，重试 %d/%d，等待 %v\n", i+1, maxRetry, delay)

		select {
		case <-time.After(delay):
			delay *= 2 // 指数退避
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("重试 %d 次后仍失败", maxRetry)
}

// 完整流水线：Audio → JSON → 中文 → TTS
// 功能：
//   - 初始化 Gemini Client
//   - 视频 → 转录 JSON
//   - JSON → 中文
//   - 中文 → MP3
func audioToJsonToTTS(ctx context.Context) error {

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		panic("GEMINI_API_KEY not set")
	}

	genaiClient, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		panic(err)
	}

	videoURL := "https://www.youtube.com/watch?v=ku-N-eS1lgM"

	// 转录
	jsonStr, err := VideoAudioToJSON(ctx, genaiClient, videoURL)
	if err != nil {
		panic(err)
	}

	// 翻译为中文
	zhText, err := TranslateToChinese(jsonStr)
	if err != nil {
		panic(err)
	}

	fmt.Println("中文内容：")
	fmt.Println(zhText)

	// 中文语音
	if err := ChineseTextToSpeechLong(ctx, zhText, "final_cn.mp3"); err != nil {
		panic(err)
	}

	fmt.Println("已生成 final_cn.mp3")

	return nil
}
