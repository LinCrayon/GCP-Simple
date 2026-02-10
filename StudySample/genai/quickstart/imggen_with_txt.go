package quickstart

import (
	"context"
	"fmt"
	"google.golang.org/genai"
	"io"
	"os"
)

/**
 * @Description
 * @Author linshengqian
 * @Date 2026/2/10 15:59
 **/
func generateMMFlashWithText(w io.Writer) error {
	// 创建一个上下文 Context 可以控制请求的生命周期，例如超时或取消。
	ctx := context.Background()

	// 创建 GenAI 客户端
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"}, // 指定 API 版本 v1
	})
	if err != nil {
		return fmt.Errorf("创建genai客户端失败: %w", err)
	}

	modelName := "gemini-2.5-flash-image"

	// 构造对话内容,使用对话式输入，每个 Content 对象可以包含多段 Part。
	contexts := []*genai.Content{
		{
			Parts: []*genai.Part{
				{
					Text: "生成一张以埃菲尔铁塔为背景，并带有烟花的图片",
				},
			},
			Role: genai.RoleUser, // 指明这条内容是用户发起的
		},
	}

	// 调用模型生成内容
	resp, err := client.Models.GenerateContent(
		ctx,
		modelName,
		contexts,
		&genai.GenerateContentConfig{
			// 指定希望返回的模态类型
			ResponseModalities: []string{
				string(genai.ModalityText),  // 文本
				string(genai.ModalityImage), // 图片
			},
			CandidateCount: int32(1), // 返回 1 个候选结果
			// 安全策略
			SafetySettings: []*genai.SafetySetting{
				{Method: genai.HarmBlockMethodProbability},               // 基于概率的安全策略
				{Category: genai.HarmCategoryDangerousContent},           // 检测危险内容
				{Threshold: genai.HarmBlockThresholdBlockMediumAndAbove}, // 阈值：中等及以上阻止
			},
		})
	if err != nil {
		return fmt.Errorf("无法生成内容: %w", err)
	}

	// 检查返回结果是否为空 ，返回的候选结果Candidates可能有多个
	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return fmt.Errorf("no candidates returned")
	}

	// 遍历返回的内容
	var fileName string
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" {
			// 如果模型返回文本，就写入 io.Writer
			fmt.Fprintln(w, part.Text)
		} else if part.InlineData != nil {
			// 如果模型返回图片（内嵌数据）
			fileName = "example-image-eiffel-tower.png"
			if err := os.WriteFile(fileName, part.InlineData.Data, 0o644); err != nil {
				return fmt.Errorf("failed to save image: %w", err)
			}
		}
	}

	// 将生成的图片文件名写入 io.Writer，方便测试或日志记录
	fmt.Fprintln(w, fileName)

	return nil
}
