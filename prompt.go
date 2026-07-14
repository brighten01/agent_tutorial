package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()
	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  os.Getenv("DEEPSEEK_API_KEY"),
		BaseURL: "https://api.deepseek.com",
		Model:   "deepseek-chat",
	})

	if err != nil {
		log.Fatal(err)
	}

	//替换对应变量
	//template := prompt.FromMessages(
	//	schema.FString,
	//	schema.SystemMessage("你是一个{role}"),
	//	schema.UserMessage("{question}"),
	//)
	//多伦对话
	template := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("你是一个{role},你的专长{exp}"),
		schema.UserMessage("{question}"),
		//加入多伦对话
		schema.AssistantMessage("我理解了让我思考一下", nil),
		schema.UserMessage("详细说明"),
	)
	//保持key同上是一致的
	val := map[string]any{
		"role":     "专业go研发",
		"exp":      "架构设计",
		"question": "设计一个简单的grpc接口服务",
	}
	messages, err := template.Format(ctx, val)
	for _, val := range messages {
		fmt.Printf("role is  %v message is %v \n", val.Role, val.Content)
	}
	content, err := chatModel.Generate(ctx, messages)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(content)
}
