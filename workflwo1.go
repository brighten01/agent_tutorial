package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/adk"
)

func main() {
	ctx := context.Background()
	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  os.Getenv("DEEPSEEK_API_KEY"),
		BaseURL: "https://api.deepseek.com",
		Model:   "deepseek-chat",
	})
	if err != nil {
		panic(err)
	}
	// 创建2个子Agent
	mainAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "MainAgent",
		Description: "负责生成初步方案",
		Instruction: "你是一个解决方案专家，请根据问题生成解决方案，如果解决方案需要改进，请说明需要改进的地方",
		Model:       chatModel,
		OutputKey:   "solution",
	})
	if err != nil {
		log.Fatalf("创建 mainAgent err is %v", err)
	}
	//agent2 批判反馈
	critiqueAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "CritiqueAgent",
		Description: "对解决方案进行批判性反馈",
		Instruction: "你是一个指令评审专家，请审查解决方案的质量，提供改进建议，如果使用方案已经足够好，请明确说明，解决方案已经完善无需改进，可以使用{solution}获取当前解决方案",
		Model:       chatModel,
		OutputKey:   "critique",
	})
	if err != nil {
		log.Fatalf("创建 critiqueAgent err is %v", err)
	}

	//创建LoopAgent
	loopAgent, err := adk.NewLoopAgent(ctx, &adk.LoopAgentConfig{
		Name:          "ReflectionAgent",
		Description:   "迭代反思智能体，通过多次迭代优化解决方案",
		SubAgents:     []adk.Agent{mainAgent, critiqueAgent},
		MaxIterations: 5,
	})

	if err != nil {
		log.Fatalf("创建 loopAgent err is %v", err)
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           loopAgent,
		EnableStreaming: false,
	})

	query := "设计一个支持高并发的网关系统"
	fmt.Printf("用户 %s \\n", query)
	iter := runner.Query(ctx, query)
	iterator := 0
	for {
		event, ok := iter.Next()
		if !ok {
			log.Fatalf("iter.Next() err is %v", err)
			break
		}

		if event.Err != nil {
			log.Fatalf("event.Err is %v", event.Err)
		}
		if event.Output != nil && event.Output.MessageOutput != nil {
			msg := event.Output.MessageOutput.Message
			if msg != nil {
				if event.AgentName == "MainAgent" {
					iterator++
					fmt.Printf("\\n迭代%d生成方案===\\n", iterator)

				} else if event.AgentName == "CritiqueAgent" {
					fmt.Printf("\\n迭代%d 批判反馈===\\n", iterator)
				}

				fmt.Printf("[%s] %s\\n", event.AgentName, msg.Content)

			}
		}
	}
}
