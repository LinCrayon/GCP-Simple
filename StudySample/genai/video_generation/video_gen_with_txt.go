package video_generation

import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/genai"
)

// generateVideoWithText 演示如何使用 Gemini Video 模型根据文本生成视频。
// 输出视频会上传到指定的 GCS URI（Google Cloud Storage）。
func generateVideoWithText(w io.Writer, outputGCSURI string) error {
	// 示例输出路径（可传入真实 GCS 路径）
	// outputGCSURI = "gs://your-bucket/your-prefix"

	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("创建 genai 客户端失败: %w", err)
	}

	// 视频生成配置
	config := &genai.GenerateVideosConfig{
		NumberOfVideos:  1,                   // 只生成一个视频
		GenerateAudio:   genai.Ptr(false),    // 关闭音频
		AspectRatio:     "16:9",              // 视频宽高比
		DurationSeconds: genai.Ptr(int32(4)), // 视频时长 4 秒
		OutputGCSURI:    outputGCSURI,        // 视频输出路径（GCS）
	}

	// 模型名称
	modelName := "veo-3.1-fast-generate-001"

	prompt := "一只猫正在读书"

	// 调用视频生成接口
	operation, err := client.Models.GenerateVideos(ctx, modelName, prompt, nil, config)
	if err != nil {
		return fmt.Errorf("启动视频生成失败: %w", err)
	}

	// 轮询直到生成操作完成
	for !operation.Done {
		time.Sleep(15 * time.Second)

		// 获取操作状态
		operation, err = client.Operations.GetVideosOperation(ctx, operation, nil)
		if err != nil {
			return fmt.Errorf("获取操作状态失败: %w", err)
		}
	}

	// 检查生成结果
	if operation.Response != nil && len(operation.Response.GeneratedVideos) > 0 {
		videoURI := operation.Response.GeneratedVideos[0].Video.URI
		fmt.Fprintln(w, videoURI)

		// 输出示例：
		// gs://your-bucket/your-prefix/videoURI
		return nil
	}

	return fmt.Errorf("视频生成失败或未返回结果")
}
