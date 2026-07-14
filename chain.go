package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// 线性执行 组建按照顺序执行
// 数据流转 自行处理问题的数据传递
// 灵活组合 支持多类型数据组合
// 类型安全 编译时类型检查

func main() {
	ctx := context.Background()
	chatTemplate := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("你是一个{role}"),
		schema.UserMessage("{question}"),
	)

	//创建chain 传递给ChatTemplate  传递给模型
	chain := compose.NewChain[map[string]any, *schema.Message]()

	chatmodel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  os.Getenv("API_KEY"),
		BaseURL: os.Getenv("BASE_URL"),
		Model:   "deepseek-chat",
	})

	chain.
		AppendChatTemplate(chatTemplate). //格式化模板
		AppendChatModel(chatmodel)        // 调用模型

	//看看是否能运行
	runnable, err := chain.Compile(ctx)

	if err != nil {
		log.Fatal("编译失败", err)
	}

	input := map[string]any{
		"role":     "go专家",
		"question": "非缓冲区channel是什么",
	}

	//模板中变量引入input
	output, err := runnable.Invoke(ctx, input)
	if err != nil {
		log.Fatal("执行失败", err)
	}
	//输出
	fmt.Println(output.Content)
	//messages, _ := chatTemplate.Format(ctx, input)
	//content, err := chatmodel.Generate(ctx, messages)
	//
	//fmt.Print(content.Content)

}
