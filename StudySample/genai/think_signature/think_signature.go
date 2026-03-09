package think_signature

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"google.golang.org/genai"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

// WeatherData 定义 Open-Meteo 返回的结构体（简化版）
type WeatherData struct {
	Current struct {
		Temperature float64 `json:"temperature_2m"` // 当前温度（单位：摄氏度）
	} `json:"current"`
}

// GeoData 用于根据城市名称获取经纬度
type GeoData struct {
	Results []struct {
		Latitude  float64 `json:"latitude"`  // 纬度
		Longitude float64 `json:"longitude"` // 经度
		Name      string  `json:"name"`      // 城市名称
	} `json:"results"`
}

// RSS 用于解析 Google News RSS
type RSS struct {
	Channel struct {
		Items []struct {
			Title string `xml:"title"` // 新闻标题
			Link  string `xml:"link"`  // 新闻链接
		} `xml:"item"`
	} `xml:"channel"`
}

func thinkSignature() {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:  "train-crayon-20260304",
		Location: "us-central1",
		Backend:  genai.BackendVertexAI,
	})
	if err != nil {
		log.Fatal("创建 genai 客户端失败： %w", err)
	}

	//////// 天气   ////////
	weatherFunc := &genai.FunctionDeclaration{
		Description: "获取指定地点的当前天气温度", //函数描述（给模型看的）
		Name:        "getWeather",
		Parameters: &genai.Schema{
			Type: "object",
			Properties: map[string]*genai.Schema{ //参数定义（JSON Schema）
				"location": {
					Type:        "string",
					Description: "城市名称，例如：伦敦",
				},
			},
			Required: []string{"location"}, // 必填字段
		},
	}

	//////// 新闻  ////////
	newsFunc := &genai.FunctionDeclaration{
		Name:        "getNews",
		Description: "获取指定主题的最新新闻",
		Parameters: &genai.Schema{
			Type: "object",
			Properties: map[string]*genai.Schema{
				"topic": {
					Type:        "string",
					Description: "新闻主题，例如: world, tech, health",
				},
			},
			Required: []string{"topic"},
		},
	}

	config := &genai.GenerateContentConfig{
		// 注册工具（Function Calling）
		Tools: []*genai.Tool{
			{
				FunctionDeclarations: []*genai.FunctionDeclaration{
					weatherFunc,
					newsFunc,
				},
			},
		},
		Temperature: genai.Ptr(float32(0.0)),
		ThinkingConfig: &genai.ThinkingConfig{
			IncludeThoughts: true, // 开启思考链
		},
	}

	prompt := "获取福建新闻和哈尔滨的天气"

	contents := []*genai.Content{
		{Parts: []*genai.Part{
			{Text: prompt},
		},
			Role: genai.RoleUser},
	}

	modelName := "gemini-2.5-flash"
	resp, err := client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		log.Fatal("failed to generate content: %w", err)
	}

	var funcCalls []*genai.FunctionCall
	for _, p := range resp.Candidates[0].Content.Parts {
		// 如果模型建议调用函数
		if p.FunctionCall != nil {
			funcCalls = append(funcCalls, p.FunctionCall)
			fmt.Printf("The model suggests to call the function %q with args: %v\n",
				p.FunctionCall.Name, p.FunctionCall.Args)
		}
	}
	if funcCalls == nil {
		log.Println("model did not suggest a function call")
	}

	var toolParts []*genai.Part
	//var weatherResult map[string]any
	//var newsResult map[string]any

	for _, fc := range funcCalls {

		var result map[string]any
		var err error

		switch fc.Name {

		case "getWeather":
			fmt.Println("正在获取天气...")
			result, err = getRealWeather(fc.Args["location"].(string))
			//weatherResult = result

		case "getNews":
			fmt.Println("正在获取新闻...")
			result, err = getNews(fc.Args["topic"].(string))
			//newsResult = result

		default:
			log.Printf("未知函数: %s\n", fc.Name)
		}

		if err != nil {
			log.Fatal(err)
		}

		// 构造 FunctionResponse
		toolParts = append(toolParts, &genai.Part{
			FunctionResponse: &genai.FunctionResponse{
				Name:     fc.Name,
				Response: result,
			},
		})
	}

	history := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{
					Text: prompt,
				},
			},
		},
		// 模型的 FunctionCall
		resp.Candidates[0].Content,
		{
			Role:  "tool",
			Parts: toolParts,
		},
	}
	//第二次调用模型生成最终回答
	resp, err = client.Models.GenerateContent(ctx, modelName, history, config)
	if err != nil {
		log.Fatal("failed to generate content: %w", err)
	}

	respText := resp.Text()
	fmt.Println(respText)

	//err = generateImage(ctx, client, weatherResult, newsResult)
	//if err != nil {
	//	log.Println("生成图片失败:", err)
	//} else {
	//	fmt.Println("图片生成成功: report.png")
	//}

}

func translate(text string) string {

	url := fmt.Sprintf(
		"https://translate.googleapis.com/translate_a/single?client=gtx&sl=auto&tl=en&dt=t&q=%s",
		url.QueryEscape(text),
	)

	resp, err := http.Get(url)
	if err != nil {
		return text
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result []any
	json.Unmarshal(body, &result)

	return result[0].([]any)[0].([]any)[0].(string)
}

// 获取真实天气的函数
func getRealWeather(location string) (map[string]any, error) {
	location = translate(location)
	// 1 获取经纬度
	geoURL := fmt.Sprintf(
		"https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1",
		url.QueryEscape(location),
	)

	resp, err := http.Get(geoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var geo GeoData
	json.Unmarshal(body, &geo)

	if len(geo.Results) == 0 {
		return nil, fmt.Errorf("找不到城市")
	}

	lat := geo.Results[0].Latitude
	lon := geo.Results[0].Longitude

	// 2 查询天气
	weatherURL := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m",
		lat, lon,
	)

	resp2, err := http.Get(weatherURL)
	if err != nil {
		return nil, err
	}
	defer resp2.Body.Close()

	body2, _ := io.ReadAll(resp2.Body)

	var data WeatherData
	json.Unmarshal(body2, &data)

	return map[string]any{
		"location":    location,
		"temperature": fmt.Sprintf("%.1f°C", data.Current.Temperature),
	}, nil
}

func getNews(topic string) (map[string]any, error) {

	topic = translate(topic)

	url := fmt.Sprintf(
		"https://news.google.com/rss/search?q=%s&hl=zh-CN&gl=CN&ceid=CN:zh-Hans",
		url.QueryEscape(topic),
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var rss RSS
	xml.Unmarshal(body, &rss)

	news := []map[string]string{}

	for i, item := range rss.Channel.Items {

		if i >= 5 {
			break
		}

		news = append(news, map[string]string{
			"title": item.Title,
			"link":  item.Link,
		})
	}

	return map[string]any{
		"topic": topic,
		"news":  news,
	}, nil
}

func generateImage(ctx context.Context,
	client *genai.Client,
	weather map[string]any,
	news map[string]any,
) error {

	location := weather["location"]
	temp := weather["temperature"]
	topic := news["topic"]

	list := news["news"].([]map[string]string)

	newsText := ""
	for i, n := range list {

		if i >= 5 {
			break
		}

		newsText += fmt.Sprintf("%d. %s\n", i+1, n["title"])
	}

	// 让 Gemini 生成海报
	prompt := fmt.Sprintf(`
设计一张科技风格的信息海报：

城市: %v
温度: %v

新闻主题: %v

新闻:
%v

要求：
- 信息图风格
- 天气图标
- 新闻列表
- 蓝色科技背景
`, location, temp, topic, newsText)

	model := "gemini-2.0-flash-preview-image-generation"

	resp, err := client.Models.GenerateContent(
		ctx,
		model,
		[]*genai.Content{
			{
				Role: "user",
				Parts: []*genai.Part{
					{Text: prompt},
				},
			},
		},
		nil,
	)

	if err != nil {
		return err
	}

	for _, part := range resp.Candidates[0].Content.Parts {

		if part.InlineData != nil {

			err := os.WriteFile(
				"report.png",
				part.InlineData.Data,
				0644,
			)

			if err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("没有返回图片")
}
