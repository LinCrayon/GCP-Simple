package embeddings

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateEmbedContentWithText 展示如何将内容嵌入文本。
func generateEmbedContentWithText(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("创建 genai 客户端失败： %w", err)
	}

	//向量长度:每条文本或文档会被编码成 3072 维的向量
	outputDimensionality := int32(3072)
	config := &genai.EmbedContentConfig{
		TaskType:             "RETRIEVAL_DOCUMENT",  //嵌入任务内容是文档检索
		Title:                "Driver's License",    //嵌入内容一个标题/标签
		OutputDimensionality: &outputDimensionality, //3072维度作为输出向量的长度传给配置
	}

	contents := []*genai.Content{
		{
			Parts: []*genai.Part{
				{
					Text: "我如何获得驾驶执照/学习许可证？",
				},
				{
					Text: "我的驾驶执照的有效期是多长时间？",
				},
				{
					Text: "驾驶员知识考试学习指南",
				},
			},
			Role: genai.RoleUser,
		},
		genai.NewContentFromText("生命的意义是什么？", genai.RoleUser), //文本字符串包装成 Content 对象
	}

	modelName := "gemini-embedding-001"
	resp, err := client.Models.EmbedContent(ctx, modelName, contents, config)
	if err != nil {
		return fmt.Errorf("无法生成内容： %w", err)
	}

	fmt.Fprintln(w, resp)

	return nil
}
