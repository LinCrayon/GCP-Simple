package translate

//
//import (
//	"context"
//	"fmt"
//	"os"
//	"strings"
//
//	translate "cloud.google.com/go/translate/apiv3"
//	"cloud.google.com/go/translate/apiv3/translatepb"
//	"google.golang.org/api/option"
//)
//
//// TranslateStandard 使用 Google Cloud Translation API（标准服务）
//func TranslateStandard(ctx context.Context, text string, sourceLang, targetLang string) (string, error) {
//	if strings.TrimSpace(text) == "" {
//		return "", fmt.Errorf("text cannot be empty")
//	}
//
//	// 设置认证
//	credsJSON, err := os.ReadFile("C:\\Users\\linshengqian\\Desktop\\GCP-Notes\\SA\\genai-sa-key.json")
//	if err != nil {
//		return "", fmt.Errorf("failed to read credentials file: %w", err)
//	}
//
//	// 使用 Translation API 客户端
//	client, err := translate.NewTranslationClient(ctx,
//		option.WithCredentialsJSON(credsJSON),
//	)
//	if err != nil {
//		return "", fmt.Errorf("failed to create translation client: %w", err)
//	}
//	defer client.Close()
//
//	// 注意：需要替换为你的项目ID
//	projectID := "train-crayon-20260104"
//	location := "global" // Translation API 通常使用 global 位置
//
//	// 构建请求
//	req := &translatepb.TranslateTextRequest{
//		Parent:             fmt.Sprintf("projects/%s/locations/%s", projectID, location),
//		SourceLanguageCode: sourceLang,
//		TargetLanguageCode: targetLang,
//		Contents:           []string{text},
//		MimeType:           "text/plain", // 纯文本
//		// 可选：指定模型类型
//		Model: "translate-llm", // 或者 "nmt"（神经网络机器翻译）
//	}
//
//	// 调用 API
//	resp, err := client.TranslateText(ctx, req)
//	if err != nil {
//		return "", fmt.Errorf("translation failed: %w", err)
//	}
//
//	// 处理响应
//	if len(resp.Translations) == 0 {
//		return "", fmt.Errorf("no translation returned")
//	}
//
//	return resp.Translations[0].TranslatedText, nil
//}
//
//// TranslateBatch 批量翻译
//func TranslateBatch(ctx context.Context, texts []string, sourceLang, targetLang string) ([]string, error) {
//	if len(texts) == 0 {
//		return []string{}, nil
//	}
//
//	// 设置认证
//	credsJSON, err := os.ReadFile("C:\\Users\\linshengqian\\Desktop\\GCP-Notes\\SA\\genai-sa-key.json")
//	if err != nil {
//		return nil, fmt.Errorf("failed to read credentials file: %w", err)
//	}
//
//	client, err := translate.NewTranslationClient(ctx,
//		option.WithCredentialsJSON(credsJSON),
//	)
//	if err != nil {
//		return nil, fmt.Errorf("failed to create translation client: %w", err)
//	}
//	defer client.Close()
//
//	projectID := "train-crayon-20260104"
//	location := "global"
//
//	req := &translatepb.TranslateTextRequest{
//		Parent:             fmt.Sprintf("projects/%s/locations/%s", projectID, location),
//		SourceLanguageCode: sourceLang,
//		TargetLanguageCode: targetLang,
//		Contents:           texts,
//		MimeType:           "text/plain",
//		Model:              "translate-llm", // 使用神经网络翻译模型
//	}
//
//	resp, err := client.TranslateText(ctx, req)
//	if err != nil {
//		return nil, fmt.Errorf("batch translation failed: %w", err)
//	}
//
//	results := make([]string, len(resp.Translations))
//	for i, translation := range resp.Translations {
//		results[i] = translation.TranslatedText
//	}
//
//	return results, nil
//}
