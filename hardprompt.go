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

type UserProfile struct {
	Name      string
	Age       int
	Interrest []string
	VIPLevel  int
}

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
		schema.SystemMessage("你是一个智能助手"),
		schema.UserMessage(`用户信息  
			"姓名 {name}
			"年龄 :{age}
			"兴趣:{interrest}
			"vip等级:{vip_level}
	根据以上推荐合适的工作行业 限制10字`),
	)
	user := &UserProfile{
		Age:       20,
		Interrest: []string{"篮球"},
		VIPLevel:  1,
		Name:      "张三",
	}

	//保持key同上是一致的
	val := map[string]any{
		"name":      user.Name,
		"age":       user.Age,
		"interrest": user.Interrest,
		"vip_level": user.VIPLevel,
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
