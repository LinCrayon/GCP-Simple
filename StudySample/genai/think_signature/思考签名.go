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
)

// WeatherData 定义 Open-Meteo 返回的结构体（简化版）
type WeatherData struct {
	Current struct {
		Temperature float64 `json:"temperature_2m"`
	} `json:"current"`
}

type GeoData struct {
	Results []struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Name      string  `json:"name"`
	} `json:"results"`
}

type RSS struct {
	Channel struct {
		Items []struct {
			Title string `xml:"title"`
			Link  string `xml:"link"`
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
		Description: "获取指定地点的当前天气温度",
		Name:        "getWeather",
		Parameters: &genai.Schema{
			Type: "object",
			Properties: map[string]*genai.Schema{
				"location": {
					Type:        "string",
					Description: "城市名称，例如：伦敦",
				},
			},
			Required: []string{"location"},
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
			IncludeThoughts: true,
		},
	}

	prompt := "中国深圳天气"

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

	var funcCall *genai.FunctionCall
	for _, p := range resp.Candidates[0].Content.Parts {
		if p.FunctionCall != nil {
			funcCall = p.FunctionCall
			fmt.Printf("The model suggests to call the function ")
			fmt.Printf("%q with args: %v\n", funcCall.Name, funcCall.Args)
		}
	}
	if funcCall == nil {
		log.Fatal("model did not suggest a function call")
	}

	var result map[string]any

	switch funcCall.Name {

	case "getWeather":
		fmt.Println("正在获取天气...")
		result, err = getRealWeather(funcCall.Args["location"].(string))

	case "getNews":
		fmt.Println("正在获取新闻...")
		result, err = getNews(funcCall.Args["topic"].(string))

	default:
		log.Fatalf("未知函数: %s", funcCall.Name)
	}

	if err != nil {
		log.Fatal(err)
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
		resp.Candidates[0].Content,
		{
			Role: "tool",
			Parts: []*genai.Part{{
				FunctionResponse: &genai.FunctionResponse{
					Name:     funcCall.Name,
					Response: result,
				},
			}},
		},
	}

	resp, err = client.Models.GenerateContent(ctx, modelName, history, config)
	if err != nil {
		log.Fatal("failed to generate content: %w", err)
	}

	respText := resp.Text()
	fmt.Println(respText)

}

func translateCity(city string) string {

	cityMap := map[string]string{
		"深圳": "Shenzhen",
		"北京": "Beijing",
		"上海": "Shanghai",
		"广州": "Guangzhou",
		"杭州": "Hangzhou",
		"东京": "Tokyo",
		"纽约": "New York",
		"伦敦": "London",
		"巴黎": "Paris",
	}

	if v, ok := cityMap[city]; ok {
		return v
	}

	return city
}

func translateTopic(topic string) string {

	m := map[string]string{
		"中国": "China",
		"美国": "USA",
		"日本": "Japan",
		"科技": "technology",
		"经济": "economy",
		"AI": "AI",
	}

	if v, ok := m[topic]; ok {
		return v
	}

	return topic
}

// 获取真实天气的函数
func getRealWeather(location string) (map[string]any, error) {
	location = translateCity(location)
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

	topic = translateTopic(topic)

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

/*funcResp := &genai.FunctionResponse{
	Name: "getWeather",
	Response: map[string]any{
		"location":         "Boston",
		"temperature":      "38",
		"temperature_unit": "F",
		"description":      "Cold and cloudy",
		"humidity":         "65",
		"wind":             `{"speed": "10", "direction": "NW"}`,
	},
}
contents = []*genai.Content{
	{
		Parts: []*genai.Part{
			{
				Text: "波士顿的天气怎么样",
			},
		},
		Role: genai.RoleUser,
	},
	{
		Parts: []*genai.Part{
			{
				FunctionCall: funcCall,
			},
		},
	},
	{
		Parts: []*genai.Part{
			{
				FunctionResponse: funcResp,
			},
		},
	},
}
*/
