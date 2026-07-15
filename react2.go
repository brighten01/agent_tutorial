package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"

	"log"
	"os"
)

//
//1.理解用户意图
//2.规划使用哪些工具
//3.执行工具调用
//4.分析工具返回结果
//5.生成最终答案

func main() {
	ctx := context.Background()
	chatMddel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		Model:   "deepseek-chat",
		BaseURL: "https://api.deepseek.com",
		APIKey:  os.Getenv("DEEPSEEK_API_KEY"),
	})

	if err != nil {
		log.Fatalf("NewChatModel err: %v", err)
	}
	timeTool := utils.NewTool(&schema.ToolInfo{
		Name:        "get_current_time",
		Desc:        "获取当前时间",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{}),
	}, func(ctx context.Context, params map[string]any) (string, error) {
		now := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("[工具执行] get_current_time -> %s\\n", now)
		return now, nil
	})
	calculator := utils.NewTool(&schema.ToolInfo{
		Name: "calculator",
		Desc: "执行简单的加减乘除",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"expression": {
				Type:     schema.String,
				Desc:     "数学表但是 例如 10+5",
				Required: true,
			},
		}),
	}, func(ctx context.Context, params map[string]any) (string, error) {
		expression := params["expression"].(string)
		result := "15"
		fmt.Printf("[工具执行] calculator(%s) ->  %s\\n", expression, result)
		return result, nil
	})
	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatMddel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{timeTool, calculator},
		},
	})
	if err != nil {
		log.Fatalf("NewAgent err: %v", err)
	}
	messages := []*schema.Message{
		schema.SystemMessage("10+1等于几"),
	}
	stream, err := agent.Stream(ctx, messages)
	defer stream.Close()

	for {
		chunk, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatalf("接受失败: %v", err)
		}
		fmt.Printf("接收到的内容: %s", chunk.Content)
	}
	if err != nil {
		log.Fatalf("Generate err: %v", err)
	}
	//fmt.Printf("最终答案: %s\n", res.Content)
}
