package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

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
