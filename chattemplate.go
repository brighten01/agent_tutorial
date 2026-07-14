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

// 思维链 大模型学习之后按照你的方式执行
func main() {
	ctx := context.Background()
	chatmodel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  os.Getenv("API_KEY"),
		BaseURL: os.Getenv("BASE_URL"),
		Model:   "deepseek-chat",
	})
	if err != nil {
		log.Fatal(err)
	}
	template := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage(`你是一个逻辑推理专家。
		1. 理解问题复述问题的要求
		2.分析：列出问题需要的步骤
		3.验证：检查答案是否合理
		4.结论：给出最终答案`),
		schema.UserMessage("{problem}"),
	)
	problem := "一个http 请求调用了一个gorutine  ,这个协程里面又启动了一个goroutine ,那么这个goroutine他们是否有关系，并且日志能通么，如何做异常的处理"
	messages, _ := template.Format(ctx, map[string]any{
		"problem": problem,
	})

	content, err := chatmodel.Generate(ctx, messages)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(content.Content)
}
