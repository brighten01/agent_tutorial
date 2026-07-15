package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/adk"
)

//  adk 创建agent  协同调用  deepseek模型

func main() {
	ctx := context.Background()
	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  os.Getenv("DEEPseek_APPID"),
		BaseURL: "https://api.deepseek.com",
		Model:   "deepseek-chat",
	})
	if err != nil {
		log.Fatalf("chatModel err is %v", err)
	}

	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "SimpleAssistant",
		Description: "一个简单的助手Agent, 能够回答问题",
		Model:       chatModel,
		ToolsConfig: adk.ToolsConfig{},
	})
	if err != nil {
		log.Fatalf("agent err is %v", err)
	}
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: false,
	})
	if err != nil {
		log.Fatalf("runner err is %v", err)
	}
	query := "浅谈非关系型数据库elasticsearch"
	fmt.Printf("用户, 查询 %s \\n", query)
	iter := runner.Query(ctx, query)
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}

		if event.Output != nil && event.Output.MessageOutput != nil {
			msg := event.Output.MessageOutput.Message
			if msg != nil {
				fmt.Println("助手:", msg.Content)
			}
		}
	}

}
