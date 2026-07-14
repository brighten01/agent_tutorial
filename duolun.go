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
	chatmodel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  os.Getenv("DEEPSEEK_API_KEY"),
		BaseURL: "https://api.deepseek.com",
		Model:   "deepseek-chat",
	})

	if err != nil {
		log.Fatal(err)
	}

	//大模型会根据你给的提示学习
	template := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("你是个情感分析助手，请分析文本的情感倾向 [正面/负面/中性] 置信度 [0/100]"),
		schema.UserMessage("这个产品我非常满意"),
		schema.AssistantMessage("情感正面 之心度95", nil),
		schema.UserMessage("一般般吧没什么特别"),
		schema.AssistantMessage("情感中性 |置信度80", nil),
		schema.UserMessage("太差了完全不能用"),
		schema.AssistantMessage("情感：负面 |置信度98", nil),
		schema.UserMessage("{text}"),
	)

	//验证这一组对话并且评估出可信度和情感色彩
	texts := []string{
		"这个框架文档不错，上手很快",
		"性能可以但是功能不多",
		"bug 多开发体验很差",
	}
	for _, text := range texts {
		messages, _ := template.Format(ctx, map[string]any{
			"text": text,
		})
		content, err := chatmodel.Generate(ctx, messages)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(content.Content)
	}

}
