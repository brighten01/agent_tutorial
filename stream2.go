package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

func main() {

	ctx := context.Background()
	chatmodel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  os.Getenv("DEEPSEEK_API_KEY"),
		BaseURL: "https://api.deepseek.com",
		Model:   "deepseek-chat",
	})
	if err != nil {
		log.Fatalf("chatmodel err is %v ", err)
	}

	searchTool := utils.NewTool(&schema.ToolInfo{
		Name: "search",
		Desc: "搜索引擎",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"query": {
				Type:     schema.String,
				Desc:     "搜索关键词",
				Required: true,
			},
		}),
	}, func(ctx context.Context, params map[string]any) (output string, err error) {
		query := params["query"]
		result := fmt.Sprintf("找到关于%s的信息 ：以下是GO控制并发理论", query)
		fmt.Printf("[工具执行] search(%s) ->  %s\\n", query, result)
		return result, nil
	})
	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatmodel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{
				searchTool,
			},
		},
	})

	if err != nil {
		log.Fatalf("agent err is %v ", err)
	}
	messages := []*schema.Message{
		schema.UserMessage("请告诉我如何控制Go协程数量及并发控制"),
	}
	fmt.Println("======用户：告诉我 如何控制协程数量=====")
	stream, err := agent.Stream(ctx, messages)
	defer stream.Close()

	for {
		chunk, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatalf("stream err is %v", err)
		}
		fmt.Print(chunk.Content)
	}

	fmt.Println("=======搜索服务完成 ======")
}
