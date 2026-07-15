package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
)

func main() {
	ctx := context.Background()

	embedder, err := openai.NewEmbedder(ctx, &openai.EmbeddingConfig{
		// 填入你的 OpenAI / 兼容OpenAI格式中转服务的 API Key
		APIKey: "ark-94d2cd4b-30a4-44d9-99a7-afa4469d9009-09252",
		// 直接填写官方模型名称，不需要ep接入点ID
		Model: "doubao-embedding-text-240715",
		// 可选：如果是国内反向代理/中转地址，替换BaseURL
		BaseURL: "https://ark.cn-beijing.volces.com/api/v3",
	})

	if err != nil {
		fmt.Printf("初始化Embedding客户端失败: %v\n", err)
		return
	}

	// 同样使用 EmbedStrings 方法，上层业务代码完全不用改动
	vectors, err := embedder.EmbedStrings(ctx, []string{"测试文本"})
	if err != nil {
		fmt.Printf("向量化调用失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 向量生成成功，向量维度：%d\n", len(vectors[0]))
}
