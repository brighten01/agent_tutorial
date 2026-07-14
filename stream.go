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
