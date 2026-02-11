package translate

import (
	"bufio"
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"io"
	"os"
	"strings"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

func translateWithReference(r io.Reader, out io.Writer) error {
	ctx := context.Background()

	//设置认证
	credsJSON, err := os.ReadFile("C:\\Users\\linshengqian\\Desktop\\GCP-Notes\\SA\\genai-sa-key.json")
	if err != nil {
		return fmt.Errorf("failed to read credentials file: %w", err)
	}
	creds, err := google.CredentialsFromJSON(ctx, credsJSON,
		"https://www.googleapis.com/auth/cloud-platform",
	)
	if err != nil {
		return fmt.Errorf("failed to load credentials: %w", err)
	}

	client, err := aiplatform.NewPredictionClient(
		ctx,
		option.WithEndpoint("us-central1-aiplatform.googleapis.com:443"),
		option.WithCredentials(creds),
	)
	if err != nil {
		return fmt.Errorf("failed to create prediction client: %w", err)
	}
	defer client.Close()

	// 读取输入
	reader := bufio.NewReader(r)
	input, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read input: %w", err)
	}
	input = strings.TrimSpace(input)

	// 4. 构造请求
	instanceMap := map[string]interface{}{
		"reference_sentence_config": map[string]interface{}{
			"reference_sentence_pair_lists": []interface{}{
				map[string]interface{}{
					"reference_sentence_pairs": []interface{}{
						map[string]interface{}{
							"source_sentence": "Deploy to production",
							"target_sentence": "部署到生产环境",
						},
					},
				},
			},
			"source_language_code": "en",
			"target_language_code": "zh",
		},
		"content": []interface{}{input},
	}

	instance, err := structpb.NewStruct(instanceMap)
	if err != nil {
		return fmt.Errorf("failed to create struct: %w", err)
	}

	req := &aiplatformpb.PredictRequest{
		Endpoint:  "projects/train-crayon-20260104/locations/us-central1/publishers/google/models/translate-llm",
		Instances: []*structpb.Value{structpb.NewStructValue(instance)},
	}

	// 调用预测 API
	resp, err := client.Predict(ctx, req)
	if err != nil {
		return fmt.Errorf("prediction failed: %w", err)
	}

	// 处理响应
	if len(resp.Predictions) > 0 {
		// 根据实际的响应结构提取翻译结果
		if prediction, ok := resp.Predictions[0].GetStructValue().AsMap()["translated_content"]; ok {
			fmt.Fprintf(out, "%v\n", prediction)
		} else {
			fmt.Fprintf(out, "%v\n", resp.Predictions)
		}
	}

	return nil
}
