package video_generation

import (
	"context"
	"fmt"
	"io"
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

	image := &genai.Image{
		GCSURI:   "gs://cloud-samples-data/generative-ai/image/flowers.png",
		MIMEType: "image/png",
	}

	config := &genai.GenerateVideosConfig{
		AspectRatio:  "16:9",
		OutputGCSURI: outputGCSURI,
	}

	modelName := "veo-3.1-fast-generate-001"
	prompt := "一簇生机勃勃的野花在阳光普照的草地上轻轻摇曳的极端特写。"
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
		videoURI := operation.Response.GeneratedVideos[0].Video.URI
		fmt.Fprintln(w, videoURI)
		return nil
	}

	// Example response:
	// gs://your-bucket/your-prefix/videoURI

	return fmt.Errorf("video generation failed or returned no results")
}
