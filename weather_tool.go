package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/schema"
)

type WeatherTool struct {
	weatherData map[string]map[string]string
}

func NewWeatherTool() *WeatherTool {
	return &WeatherTool{
		weatherData: map[string]map[string]string{
			"北京": {
				"temperature": "28°C",
				"condition":   "晴天",
				"humidity":    "45%",
				"wind":        "北风3级",
			},
			"上海": {
				"temperature": "30°C",
				"condition":   "晴天",
				"humidity":    "45%",
				"wind":        "北风3级",
			},
			"深圳": {
				"temperature": "20°C",
				"condition":   "阴天",
				"humidity":    "45%",
				"wind":        "北风3级",
			},
		},
	}
}

type WeatherParams struct {
	City string `json:"city"`
}

func (w *WeatherTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "get_weather",
		Desc: "查询指定城市的天气",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			//保持和结构体json 序列化之后名称一致
			"city": {
				Type:     schema.String,
				Desc:     "城市名称 ，例如北京、上海、深圳",
				Required: true,
			},
		}),
	}, nil
}

func (w *WeatherTool) InvokableRun(ctx context.Context, argumentJson string) (string, error) {
	var params WeatherParams
	//参数必须要json反序列化
	if err := json.Unmarshal([]byte(argumentJson), &params); err != nil {
		return "", err
	}
	weather, exists := w.weatherData[params.City]
	if !exists {
		result := map[string]string{"error": fmt.Sprintf("city %s not found", params.City)}
		resultJson, err := json.Marshal(result)
		if err != nil {
			return "", err
		}
		return string(resultJson), nil
	}
	result := map[string]string{
		"temperature": weather["temperature"],
		"condition":   weather["condition"],
		"humidity":    weather["humidity"],
		"wind":        weather["wind"],
	}
	resultJson, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	//根据情况输出
	return string(resultJson), nil
}

func main() {
	ctx := context.Background()
	weatherTool := NewWeatherTool()
	cites := []string{"北京", "上海", "广州"}
	for _, city := range cites {
		params := &WeatherParams{City: city}
		paramsJson, _ := json.Marshal(params)
		result, _ := weatherTool.InvokableRun(ctx, string(paramsJson))
		fmt.Println(result)
	}
}
