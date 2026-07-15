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
	teachAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "TechReasearcher",
		Description: "负责技术调研",
		Instruction: "你是一个研究员，请调研相关技术方案",
		Model:       chatModel,
		OutputKey:   "tech_research",
	})
	if err != nil {
		log.Fatalf("创建 TechReasearcher err is %v", err)
	}
	//agent2  市场分析
	marketAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "MarketAnalyst",
		Description: "负责市场分析",
		Instruction: "你是一个市场分析师，请分析市场趋势和竞争对手",
		Model:       chatModel,
		OutputKey:   "market_analysis",
	})

	if err != nil {
		log.Fatalf("创建 marketAgent err is %v", err)
	}

	//风险评估agent
	riskAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "RiskAssessor",
		Description: "负责风险评估",
		Instruction: "你是一个风险评估专家， 请评估项目潜在风险",
		Model:       chatModel,
		OutputKey:   "risk_assessment",
	})
	if err != nil {
		log.Fatalf("创建 riskAgent err is %v", err)
	}
	//创建ParalleAgent
	parallelAgent, err := adk.NewParallelAgent(ctx, &adk.ParallelAgentConfig{
		Name:        "ParallelAgent",
		Description: "并行执行多个子任务",
		SubAgents:   []adk.Agent{teachAgent, marketAgent, riskAgent},
	})
	if err != nil {
		log.Fatalf("创建 parallelAgent err is %v", err)
	}
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           parallelAgent,
		EnableStreaming: false,
	})
	query := "请分析下为什么阿里全面取消claud code 分析下程序猿就业市场情况及it可行性分析"
	fmt.Printf("用户 %s \\n", query)
	iter := runner.Query(ctx, query)

	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Output != nil && event.Output.MessageOutput != nil {
			msg := event.Output.MessageOutput.Message
			if msg != nil {
				fmt.Printf("[%s]: %s \\n\\n", event.AgentName, msg.Content)
			}
		}

	}
}
