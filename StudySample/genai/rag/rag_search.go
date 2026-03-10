package rag

import (
	"context"
	"fmt"
	"google.golang.org/genai"
	"log"
)

func ragSearch() {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
		Backend:     genai.BackendVertexAI,
		Project:     "train-crayon-20260304",
		Location:    "asia-northeast1",
	})

	// -------------------------------
	// 配置 RAG 检索工具
	// -------------------------------
	//RAG（Retrieval-Augmented Generation）使用的语料库 ID
	ragCorpusID := "projects/train-crayon-20260304/locations/asia-east1/ragCorpora/6917529027641081856"

	//向量检索时返回的最相似文档数量,表示在检索时取最相关的 3 条文档
	similarityK := int32(3)
	//向量相似度的阈值 ,只有相似度大于等于 0.5 的文档才会被认为相关
	vectorThreshold := 0.5
	//开启上下文存储
	//storeContext := true

	ragTool := &genai.Tool{
		Retrieval: &genai.Retrieval{ // 指定这是一个检索工具
			VertexRAGStore: &genai.VertexRAGStore{ // 使用 Vertex AI RAG 存储
				RAGResources: []*genai.VertexRAGStoreRAGResource{ // 关联一个或多个 RAG 语料库资源
					{
						RAGCorpus: ragCorpusID,
					},
				},
				// 配置检索行为
				RAGRetrievalConfig: &genai.RAGRetrievalConfig{
					TopK: &similarityK,
				},
				//StoreContext:            &storeContext,
				SimilarityTopK:          &similarityK,     // 额外指定相似度排序的 TopK 数量（一般与 TopK 相同）
				VectorDistanceThreshold: &vectorThreshold, // 设定向量相似度的阈值
			},
		},
	}

	// -------------------------------
	// 配置生成内容参数
	// -------------------------------

	config := genai.GenerateContentConfig{
		ResponseModalities: []string{
			string(genai.ModalityText),
		},
		Tools: []*genai.Tool{
			//{
			//	GoogleSearch: &genai.GoogleSearch{},
			//},
			ragTool, // <-- RAG 检索工具
		},
	}

	// -------------------------------
	// 调用生成接口
	// -------------------------------
	prompt := "RAG中有个文件istio.md,里面写了什么内容"
	response, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		&config,
	)
	if err != nil {
		log.Fatalf("生成内容失败: %v", err)
	}

	// -------------------------------
	// 输出结果
	// -------------------------------
	if len(response.Candidates) > 0 && response.Candidates[0].Content != nil {
		for i, part := range response.Candidates[0].Content.Parts {
			if part.Text != "" {
				fmt.Printf("Part %d: %s\n", i+1, part.Text)
			}
		}
	} else {
		fmt.Println("没有生成内容")
	}
}
