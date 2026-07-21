package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

/*
入参：JSON 字符串 → 自动反序列化 → 你的 *TimeParams 结构体，传给 GetCurrentTime
返回：你的 *TimeResult 结构体 → 自动序列化 → JSON 字符串，作为 outputJson 返回

utils.NewTool 内部做了什么（底层封装逻辑）：
通过反射读取你的业务函数：
第一个入参结构体：TimeParams（入参模型）
返回第一个结构体：TimeResult（出参模型）
自动生成 Info() 方法：直接复用你传入的 schema.ToolInfo；
自动生成完整 InvokableRun 胶水逻辑（和你手动写的逻辑一模一样，只是封装在框架内部）：
*/
type TimeParams struct {
	Format string `json:"format"`
}
type TimeResult struct {
	CurrentTime string `json:"current_time"`
}

func GetCurrentTime(ctx context.Context, params *TimeParams) (*TimeResult, error) {
	now := time.Now()
	var result string
	switch params.Format {
	case "date":
		result = now.Format("2006-01-02")
	case "time":
		result = now.Format("15:04:05")
	default:
		result = now.Format("2006-01-02 15:04:05")

	}
	return &TimeResult{CurrentTime: result}, nil
}

func main() {
	ctx := context.Background()
	timeTool := utils.NewTool(&schema.ToolInfo{
		Name: "get_current_time",
		Desc: "获取当前时间",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"format": {
				Type:     schema.String,
				Desc:     "时间格式 date(日期), time(时间), datetime(日期时间)",
				Required: false,
			},
		}),
	}, GetCurrentTime)

	timeFormats := []string{"date", "time", "datetime"}
	for _, format := range timeFormats {
		params := &TimeParams{Format: format}
		b, _ := json.Marshal(params)
		//参数必须是字符串 LLM 做的内部转换 照做就行了
		outputJson, err := timeTool.InvokableRun(ctx, string(b))
		if err != nil {

		}
		fmt.Println(string(outputJson))
	}
}
