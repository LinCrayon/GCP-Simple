package video_generation

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"google.golang.org/genai"
)

// generateVideoFromImage 从图像生成视频
func generateVideoFromImage(w io.Writer, outputGCSURI string) error {
	//outputGCSURI = "gs://your-bucket/your-prefix"
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	img, err := os.ReadFile("fruit.png")
	if err != nil {
		return fmt.Errorf("读取输入图片失败: %w", err)
	}

	image := &genai.Image{
		ImageBytes: img,
		MIMEType:   "image/png",
	}

	config := &genai.GenerateVideosConfig{
		NumberOfVideos:  1,                   // 只生成一个视频
		GenerateAudio:   genai.Ptr(false),    // 关闭音频
		AspectRatio:     "16:9",              // 视频宽高比
		DurationSeconds: genai.Ptr(int32(4)), // 视频时长 4 秒
	}

	modelName := "veo-3.1-fast-generate-001"
	prompt := "一段电影级特写视频：木质果盘中装着红色苹果和紫色葡萄，放在明亮的厨房台面上。\n柔和的自然阳光缓慢变化，形成细腻的光影过渡。\n镜头缓慢推进，背景虚化。\n葡萄轻微晃动，仿佛被微风拂过。\n写实风格，高细节，画面稳定流畅。。"
	operation, err := client.Models.GenerateVideos(ctx, modelName, prompt, image, config)
	if err != nil {
		return fmt.Errorf("failed to start video generation: %w", err)
	}

	// Polling until the operation is done
	for !operation.Done {
		time.Sleep(15 * time.Second)
		operation, err = client.Operations.GetVideosOperation(ctx, operation, nil)
		if err != nil {
			return fmt.Errorf("failed to get operation status: %w", err)
		}
	}

	if operation.Response != nil && len(operation.Response.GeneratedVideos) > 0 {
		//videoURI := operation.Response.GeneratedVideos[0].Video.URI
		videoBytes := operation.Response.GeneratedVideos[0].Video.VideoBytes
		outputFile := "output-video.mp4"
		if err := os.WriteFile(outputFile, videoBytes, 0644); err != nil {
			return fmt.Errorf("failed to write video file: %w", err)
		}

		fmt.Fprintf(
			w,
			"Video saved locally: %s (%d bytes)\n",
			outputFile,
			len(videoBytes),
		)
		return nil
	}

	// Example response:
	// gs://your-bucket/your-prefix/videoURI

	return fmt.Errorf("video generation failed or returned no results")
}
