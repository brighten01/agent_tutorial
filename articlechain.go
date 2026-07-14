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

// 优化文章 工作流 1。总结文章大纲 2.扩写文章 3.文章润色 4。格式化输出
type ArticleRequest struct {
	Topic    string
	Keywords []string
	Length   int
}

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
	chain := compose.NewChain[ArticleRequest, string]()
	chain.AppendLambda(compose.InvokableLambda(func(ctx context.Context, req ArticleRequest) (string, error) {
		fmt.Println("===生成文章大纲====")
		template := prompt.FromMessages(
			schema.FString,
			schema.SystemMessage("你是一个专业的内容策划师。请根据主题生成文章大纲"),
			schema.UserMessage("主题{topic}\\n 关键词{keywords}\\n 氢生成一个3-5点的文章大纲。"),
		)
		message, _ := template.Format(ctx, map[string]any{
			"topic":    req.Topic,
			"keywords": fmt.Sprintf("%v", req.Keywords),
		})
		res, err := chatmodel.Generate(ctx, message)
		if err != nil {
			log.Fatalf("chatmodel err is %v ", err)
		}
		fmt.Println(" 大纲", res.Content)
		return res.Content, nil
	})).
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, outline string) (string, error) {
			fmt.Println("===扩写内容====")
			template := prompt.FromMessages(
				schema.FString,
				schema.SystemMessage("你是一个专业作者。请根据文章大纲生成完整文章，要求内容清晰，逻辑完整"),
				schema.UserMessage("大纲{outline}\\n 请扩写一片200字左右的文章"),
			)

			message, _ := template.Format(ctx, map[string]any{
				"outline": outline,
			})
			res, err := chatmodel.Generate(ctx, message)
			if err != nil {
				log.Fatalf("chatmodel err is %v ", err)
			}
			fmt.Println(" 文章内容", res.Content)
			return res.Content, nil
			//润色优化
		})).AppendLambda(compose.InvokableLambda(func(ctx context.Context, draft string) (string, error) {
		fmt.Println("===润色优化====")
		template := prompt.FromMessages(
			schema.FString,
			schema.SystemMessage("你是一个专业的编辑，优化文章表达，让文章通畅自然、生动"),
			schema.UserMessage("文章{draft}\\n 请尽兴润色"),
		)

		message, _ := template.Format(ctx, map[string]any{
			"draft": draft,
		})
		res, err := chatmodel.Generate(ctx, message)
		if err != nil {
			log.Fatalf("chatmodel err is %v ", err)
		}
		fmt.Println("润色完成")
		return res.Content, nil
		//格式化输出
	})).AppendLambda(compose.InvokableLambda(func(ctx context.Context, article string) (string, error) {
		fmt.Println("===格式化输出====")
		formated := fmt.Sprintf("#生成的文章\\n\\n%s \\n\\n----\\n ** 又框架生成", article)
		return formated, nil
	}))
	runable, err := chain.Compile(ctx)

	if err != nil {
		log.Fatalf("runable err is %v ", err)
	}
	request := ArticleRequest{
		Topic:    "人工智能在软件开发中的作用",
		Keywords: []string{"chatgpt", "ai"},
		Length:   200,
	}
	content, err := runable.Invoke(ctx, request)
	if err != nil {
		log.Fatalf("invoke err is %v ", err)
	}
	fmt.Println(content)

}
