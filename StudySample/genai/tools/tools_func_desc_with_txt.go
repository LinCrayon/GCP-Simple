package tools

import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateWithFuncCall 函数调用：演示 Gemini 模型如何不是直接回答问题，而是先“建议调用一个函数”，再利用函数返回的数据生成最终回答。
// 模型会根据用户问题判断是否需要调用某个函数，并给出函数名和参数。
// 本示例使用“获取天气”的虚拟函数，并用模拟数据完成一次完整对话。
func generateWithFuncCall(w io.Writer) error {

	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("创建 genai 客户端失败: %w", err)
	}

	// 定义一个可供模型调用的函数：获取指定城市的当前天气
	weatherFunc := &genai.FunctionDeclaration{

		Description: "返回指定地点的当前天气信息",
		// 函数名称（模型会在 FunctionCall 中返回该名称）
		Name: "getCurrentWeather",
		// 函数参数结构定义（JSON Schema）
		Parameters: &genai.Schema{
			Type: "object",
			Properties: map[string]*genai.Schema{
				// location 表示城市名称
				"location": {Type: "string"},
			},
			// location 为必填参数
			Required: []string{"location"},
		},
	}

	// 生成配置
	config := &genai.GenerateContentConfig{
		// 向模型注册可调用的函数
		Tools: []*genai.Tool{
			{
				FunctionDeclarations: []*genai.FunctionDeclaration{
					weatherFunc,
				},
			},
		},
		// 使用较低温度，保证输出稳定、可预测
		Temperature: genai.Ptr(float32(0.0)),
	}

	modelName := "gemini-2.5-flash"

	// 第一轮对话：用户提问
	contents := []*genai.Content{
		{
			Role: genai.RoleUser,
			Parts: []*genai.Part{
				{Text: "波士顿现在的天气怎么样？"},
			},
		},
	}

	// 第一次调用模型：让模型决定是否需要调用函数
	resp, err := client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		return fmt.Errorf("生成内容失败: %w", err)
	}

	// 从模型返回内容中查找 FunctionCall
	var funcCall *genai.FunctionCall
	for _, p := range resp.Candidates[0].Content.Parts {
		if p.FunctionCall != nil {
			funcCall = p.FunctionCall

			// 输出模型建议调用的函数及参数
			fmt.Fprint(w, "模型建议调用函数 ")
			fmt.Fprintf(w, "%q，参数为：%v\n", funcCall.Name, funcCall.Args)

			// 示例输出：
			// 模型建议调用函数 "getCurrentWeather"，参数为：map[location:Boston]
		}
	}

	// 如果模型未建议任何函数调用，则返回错误
	if funcCall == nil {
		return fmt.Errorf("模型未建议调用任何函数")
	}

	// ------------------------------
	// 模拟外部 API 的返回结果
	// ------------------------------
	// 在真实业务中，这里应该调用真实的天气 API
	// 本示例使用模拟数据（Synthetic Data）
	funcResp := &genai.FunctionResponse{
		// 对应函数名称
		Name: "getCurrentWeather",
		// 函数返回结果（任意 JSON 结构）
		Response: map[string]any{
			"location":         "Boston",
			"temperature":      "38",
			"temperature_unit": "F",
			"description":      "寒冷且多云",
			"humidity":         "65",
			"wind":             `{"speed": "10", "direction": "西北风"}`,
		},
	}

	// 将完整的对话历史 + 函数调用结果传回模型
	contents = []*genai.Content{
		{
			Role: genai.RoleUser,
			Parts: []*genai.Part{
				{Text: "波士顿现在的天气怎么样？"},
			},
		},
		{
			Parts: []*genai.Part{
				{FunctionCall: funcCall},
			},
		},
		{
			Parts: []*genai.Part{
				{FunctionResponse: funcResp},
			},
		},
	}

	// 第二次调用模型：让模型基于函数返回结果生成自然语言回复
	resp, err = client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		return fmt.Errorf("生成最终回复失败: %w", err)
	}

	// 获取模型生成的最终文本结果
	respText := resp.Text()

	// 输出最终回答
	fmt.Fprintln(w, respText)

	// 示例输出：
	// 波士顿目前天气寒冷且多云，气温约为 38 华氏度，湿度为 65%，
	// 西北风，风速约为 10。
	return nil
}
