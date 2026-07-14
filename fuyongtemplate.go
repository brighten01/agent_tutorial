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

type PromptTemplate struct {
}

func (p *PromptTemplate) Translator(sourcelang, targetlang string) prompt.ChatTemplate {
	return prompt.FromMessages(
		schema.FString,
		schema.SystemMessage(fmt.Sprintf("你是一个专业的翻译助手。请将%s翻译成%s。\\n"+
			"要求:\\n"+
			"1.保持原文的语气和风格\\n"+
			"2.确保翻译准确、流畅\\n"+
			"3.只返回翻译结果，不要添加解释", sourcelang, targetlang)),
		schema.UserMessage("{text}"),
	)
}
func (p *PromptTemplate) CodeReviewer(lang string) prompt.ChatTemplate {
	return prompt.FromMessages(schema.FString,
		schema.UserMessage(fmt.Sprintf("你是一个资深的%s开发专家。请审查以下代码，并提供:\\n"+
			"1.潜在的bug或问题\\n"+
			"2.性能优化建议\\n"+
			"3.代码风格改进建议\n"+
			"4.安全性评估", lang)),
		schema.UserMessage("请审查以下代码: \\n\\n{language}\\n{code}\\n"),
	)
}

// 技术面试官模板
func (p *PromptTemplate) TechInterviewer(position, level string) prompt.ChatTemplate {
	return prompt.FromMessages(
		schema.FString,
		schema.SystemMessage(fmt.Sprintf("你是一位%s职位的面试官，针对%s级别的候选人。\\n"+
			"请根据候选人的回答:\\n"+
			"1.评估答案的准确性和深度\n"+
			"2.提出有针对性的追问\n"+
			"3.给出建设性的反馈", position, level)),
		schema.UserMessage("候选人回答:{answer}\\n\\n请评估并追问。"),
	)
}

func main() {
	ctx := context.Background()
	templates := &PromptTemplate{}
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

	fmt.Println("===翻译测试====")
	translateTemplate := templates.Translator("中文", "英文")
	messages, _ := translateTemplate.Format(ctx, map[string]any{
		"text": "一个强大的框架eion",
	})
	content, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("翻译结果 ", content.Content)

	fmt.Println("====代码审计测试===")
	codeTemplate := templates.CodeReviewer("go")
	messages, _ = codeTemplate.Format(ctx, map[string]any{
		"language": "go",
		"code":     "package main\n\nfunc main (){ fmt.Println(111)}",
	})

	content, _ = chatModel.Generate(ctx, messages)
	fmt.Println(" 代码设计结果", content.Content)

}
