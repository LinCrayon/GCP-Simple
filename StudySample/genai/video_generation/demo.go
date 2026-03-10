package video_generation

import (
	"context"
	"fmt"
	"google.golang.org/genai"
	"log"
)

func main() {
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
	ragCorpusID := "projects/train-crayon-20260304/locations/asia-east1/ragCorpora/6917529027641081856"

	similarityK := int32(3)
	vectorThreshold := 0.5

	ragTool := &genai.Tool{
		Retrieval: &genai.Retrieval{
			VertexRAGStore: &genai.VertexRAGStore{
				RAGResources: []*genai.VertexRAGStoreRAGResource{
					{
						RAGCorpus: ragCorpusID,
					},
				},
				RAGRetrievalConfig: &genai.RAGRetrievalConfig{
					TopK: &similarityK,
				},
				SimilarityTopK:          &similarityK,
				VectorDistanceThreshold: &vectorThreshold,
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
	response, err := client.Models.GenerateContent(ctx,
		"gemini-2.5-flash",
		genai.Text(prompt), &config,
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
