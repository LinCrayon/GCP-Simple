package image_generation

import (
	"context"
	"fmt"
	"google.golang.org/genai"
	"io"
	"os"
)

/**
 * @Description 编辑图像
 * @Author linshengqian
 **/
//通过“文本 + 输入图片”的方式对图片进行编辑（如风格化、卡通化等），并将生成后的图片保存为本地文件。
func generateImageMMFlashEditWithTextImg(w io.Writer) error {
	// TODO(developer): Update below lines
	outputFile := "bw-example-image.png"
	inputFile := "example-image-eiffel-tower.png"
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("创建 genai 客户端失败: %w", err)
	}

	// 读取本地输入图片（作为编辑参考图）
	image, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("读取输入图片失败: %w", err)
	}

	// 使用支持图像编辑的多模态模型
	modelName := "gemini-2.5-flash-image"
	// 文本提示：描述你希望模型如何编辑这张图片
	prompt := "请将这张图片编辑成卡通风格。"
	contents := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{Text: prompt}, //文本指令
				{InlineData: &genai.Blob{ //输入图片
					MIMEType: "image/png",
					Data:     image,
				}},
			},
		},
	}

	//执行图像编辑
	resp, err := client.Models.GenerateContent(ctx,
		modelName,
		contents,
		&genai.GenerateContentConfig{
			// 同时请求返回文本说明 + 编辑后的图片
			ResponseModalities: []string{
				string(genai.ModalityText),
				string(genai.ModalityImage),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("生成内容失败: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return fmt.Errorf("模型未生成任何内容")
	}

	// 遍历模型返回的内容片段
	for _, part := range resp.Candidates[0].Content.Parts {
		// 如果是文本内容（通常是编辑说明）
		if part.Text != "" {
			fmt.Fprintln(w, part.Text)
			// 如果是生成的图片数据
		} else if part.InlineData != nil {
			if len(part.InlineData.Data) > 0 {
				// 将生成的图片保存为文件
				if err := os.WriteFile(outputFile, part.InlineData.Data, 0644); err != nil {
					return fmt.Errorf("failed to save image: %w", err)
				}
				fmt.Fprintln(w, outputFile)
			}
		}
	}

	// 示例输出（文本部分）：
	// 这是为你生成的卡通风格埃菲尔铁塔图片！
	// 编辑说明：
	// - 使用更粗的线条简化了铁塔结构
	// - 提高了颜色饱和度，使画面更具卡通感
	// - 弱化了真实光影，增强插画风格
	return nil
}
