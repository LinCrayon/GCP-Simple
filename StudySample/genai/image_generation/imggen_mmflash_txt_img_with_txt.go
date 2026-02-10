package image_generation

import (
	"context"
	"fmt"
	"google.golang.org/genai"
	"io"
	"os"
	"path/filepath"
)

/**
 * @Description 生成交错图像和文本, 例如，您可以为生成的食谱的每个步骤生成相应的图像，而无需向模型发出单独的请求。
 * @Author linshengqian
 **/
//根据一段文本提示同时生成【文字 + 图片】，
func generateMMFlashTxtImgWithText(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("创建 genai 客户端失败: %w", err)
	}

	// 使用支持“文本 + 图像”多模态输出的模型
	modelName := "gemini-2.5-flash-image"
	//构造请求内容（Prompt）
	contents := []*genai.Content{
		{
			Parts: []*genai.Part{
				{
					Text: "生成一份带插图的西班牙海鲜饭（Paella）食谱。" +
						"在生成文字步骤的同时，为每个步骤配上示意图片。",
				},
			},
			Role: genai.RoleUser,
		},
	}
	//生成内容
	resp, err := client.Models.GenerateContent(ctx,
		modelName,
		contents,
		&genai.GenerateContentConfig{
			//	指定返回内容的模态类型：
			ResponseModalities: []string{
				string(genai.ModalityText),
				string(genai.ModalityImage),
			},
			CandidateCount: int32(1), //一个候选结果
		},
	)
	if err != nil {
		return fmt.Errorf("生成内容失败: %w", err)
	}

	// 校验返回结果是否有效
	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return fmt.Errorf("未返回任何候选内容")
	}

	// 输出目录（当前目录）
	outputFolder := ""

	// 创建 Markdown 文件，用于保存文字 + 图片引用
	mdFile := filepath.Join(outputFolder, "paella-recipe.md")
	fp, err := os.Create(mdFile)
	if err != nil {
		return fmt.Errorf("创建 Markdown 文件失败: %w", err)
	}
	defer fp.Close()

	// 遍历模型返回的内容片段（Parts）
	for i, part := range resp.Candidates[0].Content.Parts {
		// 如果是文本内容，直接写入 Markdown
		if part.Text != "" {
			if _, err := fp.WriteString(part.Text); err != nil {
				return fmt.Errorf("写入文本失败: %w", err)
			}
			// 如果是图片内容（InlineData）
		} else if part.InlineData != nil {
			// 生成图片文件名
			imgFile := filepath.Join(outputFolder, fmt.Sprintf("example-image-%d.png", i+1))
			// 将图片字节数据写入 PNG 文件
			if err := os.WriteFile(imgFile, part.InlineData.Data, 0644); err != nil {
				return fmt.Errorf("failed to save image: %w", err)
			}
			// 在 Markdown 中插入图片引用
			if _, err := fp.WriteString(fmt.Sprintf("![image](%s)", filepath.Base(imgFile))); err != nil {
				return fmt.Errorf("写入图片引用失败: %w", err)
			}
		}
	}

	// 将生成的 Markdown 文件路径输出到 writer
	fmt.Fprintln(w, mdFile)

	// 示例效果：
	// - 生成一个 paella-recipe.md 文件
	// - 文件中包含详细的烹饪步骤说明
	// - 每个步骤旁边配有模型生成的示意图片
	return nil
}
