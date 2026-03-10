package rag

//
//import (
//	"context"
//	"fmt"
//	"log"
//
//	"google.golang.org/genai"
//	"google.golang.org/genai/vertexai/preview/rag"
//)
//
//func main() {
//	ctx := context.Background()
//
//	// -------------------------------
//	// 1️⃣ 初始化 Vertex AI 客户端
//	// -------------------------------
//	client, err := genai.NewClient(ctx, &genai.ClientConfig{
//		Backend:  genai.BackendVertexAI,
//		Project:  "train-crayon-20260304",
//		Location: "asia-northeast1",
//	})
//	if err != nil {
//		log.Fatalf("创建 Vertex AI client 失败: %v", err)
//	}
//
//	// -------------------------------
//	// 2️⃣ 定义要导入的文件
//	// 支持 Google Drive 或 GCS 链接
//	// -------------------------------
//	paths := []string{
//		"https://drive.google.com/file/d/YOUR_FILE_ID/view?usp=sharing",
//		"gs://my_bucket/my_files_dir/",
//	}
//
//	// -------------------------------
//	// 3️⃣ 创建 RAG Corpus
//	// -------------------------------
//	displayName := "my_rag_corpus"
//
//	genai.VertexRAGStore{
//		RAGCorpora:              nil,
//		RAGResources:            nil,
//		RAGRetrievalConfig:      nil,
//		SimilarityTopK:          nil,
//		StoreContext:            nil,
//		VectorDistanceThreshold: nil,
//	}
//
//
//	embeddingModelConfig := genai.VertexRAGStore{}.EmbeddingModelConfig{
//		PublisherModel: "publishers/google/models/text-embedding-005",
//	}
//
//	ragCorpus, err := rag.CreateCorpus(ctx, rag.CreateCorpusConfig{
//		DisplayName:          displayName,
//		EmbeddingModelConfig: embeddingModelConfig,
//	})
//	if err != nil {
//		log.Fatalf("创建 RagCorpus 失败: %v", err)
//	}
//
//	fmt.Println("创建 RagCorpus 成功:", ragCorpus.Name)
//
//	// -------------------------------
//	// 4️⃣ 导入文件到 RAG Corpus
//	// -------------------------------
//	transformationConfig := rag.TransformationConfig{
//		ChunkingConfig: rag.ChunkingConfig{
//			ChunkSize:    512,
//			ChunkOverlap: 100,
//		},
//	}
//
//	if err := rag.ImportFiles(ctx, ragCorpus.Name, paths, &transformationConfig); err != nil {
//		log.Fatalf("导入文件失败: %v", err)
//	}
//
//	fmt.Println("文件导入完成")
//
//	// -------------------------------
//	// 5️⃣ 配置 RAG 检索工具
//	// -------------------------------
//	similarityK := int32(3)
//	vectorThreshold := 0.5
//
//	ragTool := &genai.Tool{
//		Retrieval: &genai.Retrieval{
//			VertexRAGStore: &genai.VertexRagStore{
//				RAGResources: []*genai.VertexRagStoreRAGResource{
//					{RAGCorpus: ragCorpus.Name},
//				},
//				SimilarityTopK:          &similarityK,
//				VectorDistanceThreshold: &vectorThreshold,
//			},
//		},
//	}
//
//	// -------------------------------
//	// 6️⃣ 使用 Gemini 生成回答
//	// -------------------------------
//	config := genai.GenerateContentConfig{
//		ResponseModalities: []string{
//			string(genai.ModalityText),
//		},
//		Tools: []*genai.Tool{ragTool},
//	}
//
//	prompt := "What is RAG and why is it helpful?"
//
//	response, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), &config)
//	if err != nil {
//		log.Fatalf("生成内容失败: %v", err)
//	}
//
//	// 输出结果
//	fmt.Println("=== 生成文本 ===")
//	for i, part := range response.Candidates[0].Content.Parts {
//		fmt.Printf("Part %d: %s\n", i+1, part.Text)
//	}
//}
