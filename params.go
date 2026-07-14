package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/schema"
)

func main() {

	fmt.Println("======基础配置=====")
	BasicChat()
	fmt.Println("===== 高级配置=====")
	AdvancedChat()
	fmt.Println("==小说作家==")
	creativeExample()
}

func BasicChat() {
	ctx := context.Background()
	chatmodel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  os.Getenv("DEEPSEEK_API_KEY"),
		BaseURL: "https://api.deepseek.com",
		Model:   "deepseek-chat",
	})

	if err != nil {
		log.Fatalf("chatmodel err is %v ", err)
	}

	messages := []*schema.Message{
		schema.SystemMessage("你是一个技术博主"),
		schema.UserMessage("请你介绍一下go channel相关引用控制在100字以内"),
	}
	stream, err := chatmodel.Stream(ctx, messages)
	if err != nil {
		log.Fatalf("stream  output err is %v ", err)
	}
	//尽量腾出缓冲区
	defer stream.Close()

	var fullContent strings.Builder
	for {
		streamContent, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				//输出结束
				break
			}
			log.Fatalf("strea recv err is %v ", err)
		}
		//注释内容流式输出
		//fmt.Print(streamContent.Content)
		fullContent.WriteString(streamContent.Content)
	}

	fmt.Println("完整输出", fullContent.String())
}

func AdvancedChat() {
	ctx := context.Background()
	chatmodel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:      os.Getenv("DEEPSEEK_API_KEY"),
		BaseURL:     "https://api.deepseek.com",
		Model:       "deepseek-chat",
		Temperature: 0.7, // 控制输出随机性，范围 [0.0, 2.0]，越高越随机
		TopP:        0.9, // 核采样参数，范围 [0.0, 1.0]，越低越聚焦
		MaxTokens:   500, // 限制最大生成 token 数量，范围 [1, 8192]
	})
	if err != nil {
		log.Fatalf("chatmodel err is %v ", err)
	}
	messages := []*schema.Message{
		schema.SystemMessage("你是一个技术博主"),
		schema.UserMessage("请你介绍一下go map相关引用控制在100字以内"),
	}
	content, err := chatmodel.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("stream  output err is %v ", err)
	}
	fmt.Println(content)
	printTokenUsage(content)

}

func creativeExample() {
	ctx := context.Background()
	chatmodel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:           os.Getenv("DEEPSEEK_API_KEY"),
		BaseURL:          "https://api.deepseek.com",
		Model:            "deepseek-chat",
		Temperature:      0.7, // 控制输出随机性，范围 [0.0, 2.0]，越高越随机
		TopP:             0.9, // 核采样参数，范围 [0.0, 1.0]，越低越聚焦
		MaxTokens:        500, // 限制最大生成 token 数量，范围 [1, 8192]
		PresencePenalty:  0.3,
		FrequencyPenalty: 0.3,
	})

	if err != nil {
		log.Fatalf("chatmodel err is %v ", err)
	}
	messages := []*schema.Message{
		schema.SystemMessage("你是一个技术博主"),
		schema.UserMessage("写一个架构师成长故事 限制字数100字"),
	}
	content, err := chatmodel.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("generate content err is %v ", err)
	}

	fmt.Println(content.Content)
	printTokenUsage(content)

}
func printTokenUsage(content *schema.Message) {
	if content.ResponseMeta != nil && content.ResponseMeta.Usage != nil {
		fmt.Printf("\n nToken 使用统计:\\n")
		fmt.Printf(" 输入tokens  %d \\n", content.ResponseMeta.Usage.PromptTokens)
		fmt.Printf(" 输出tokens %d \\n", content.ResponseMeta.Usage.CompletionTokens)
		fmt.Println(" 总计 tokens ", content.ResponseMeta.Usage.TotalTokens)
		if content.ResponseMeta.Usage.PromptTokenDetails.CachedTokens > 0 {
			fmt.Printf("  缓存 Token: %d\\n", content.ResponseMeta.Usage.PromptTokenDetails.CachedTokens)

		}
	}
}
