package tools

import (
	"context"
	"fmt"
	"google.golang.org/genai"
	"io"
	"os"
	"strings"
)

/**
 * @Description
 * @Author linshengqian
 * @Date 2026/3/5 9:46
 **/
func generateWithGoogleMap(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		return fmt.Errorf("创建 genai 客户端失败: %w", err)
	}

	prompt := "距离这里 15 分钟步行范围内有哪些最好的意大利餐厅？"

	contents := []*genai.Content{
		{
			Role: genai.RoleUser,
			Parts: []*genai.Part{
				{Text: prompt},
			},
		},
	}
	Latitude := 34.050481
	Longitude := -118.248526

	lat := &Latitude
	long := &Longitude

	// 生成配置： 启用 GoogleMaps 工具
	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{
				GoogleMaps: &genai.GoogleMaps{}, //启用 Google Maps 功能
			},
		},
		ToolConfig: &genai.ToolConfig{
			RetrievalConfig: &genai.RetrievalConfig{ //用来指定检索相关的参数，比如位置
				LatLng: &genai.LatLng{
					Latitude:  lat,  //维度
					Longitude: long, //经度
				},
			},
		},
		// 温度设为 0，保证推理与代码生成稳定
		Temperature: genai.Ptr(float32(0.0)),
	}

	modelName := "gemini-2.5-flash"

	// 调用模型生成内容
	resp, err := client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		return fmt.Errorf("生成内容失败: %w", err)
	}

	// 输出文本
	fmt.Println("Generated Response:")
	fmt.Println(resp.Candidates[0].Content.Parts[0].Text)

	// 输出 grounding 来源
	if grounding := resp.Candidates[0].GroundingMetadata; grounding != nil { //信息来源
		if len(grounding.GroundingChunks) > 0 {
			fmt.Fprintln(w, strings.Repeat("-", 40))
			fmt.Fprintln(w, "Sources:")
			for _, chunk := range grounding.GroundingChunks {
				if chunk.Maps != nil {
					fmt.Fprintf(w, "- [%s](%s)\n", chunk.Maps.Title, chunk.Maps.URI)
				}
			}
		}
	}

	return nil
}
