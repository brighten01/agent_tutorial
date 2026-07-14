package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/compose"
)

func main() {
	ctx := context.Background()
	branchCondtion := func(ctx context.Context, input map[string]any) (string, error) {
		language := input["language"].(string)
		language = strings.ToLower(language)
		if language == "go" || language == "golang" {
			return "go_branch", nil

		}
		if language == "python" {
			return "python_branch", nil
		}
		return "other_lang_branch", nil
	}
	//Go分之处理
	goBranch := compose.InvokableLambda(func(ctx context.Context, input map[string]any) (output map[string]any, err error) {
		input["advice"] = "推荐框架进行开发"
		input["features"] = []string{"协程", "channel"}
		return input, nil
	})

	pythonBranch := compose.InvokableLambda(func(ctx context.Context, input map[string]any) (output map[string]any, err error) {
		input["advice"] = "推荐框架lanchain进行开发"
		input["features"] = []string{"生态丰富", "社区活跃"}
		return input, nil
	})

	otherBranch := compose.InvokableLambda(func(ctx context.Context, input map[string]any) (output map[string]any, err error) {
		return input, nil
	})

	chain := compose.NewChain[map[string]any, map[string]any]()

	chain.AppendLambda(compose.InvokableLambda(func(ctx context.Context, input map[string]any) (output map[string]any, err error) {
		fmt.Println("====开始处理")
		return input, nil
	})).
		AppendBranch(compose.NewChainBranch(branchCondtion).
			AddLambda("go_branch", goBranch).
			AddLambda("python_branch", pythonBranch).
			AddLambda("other_branch", otherBranch),
		).AppendLambda(compose.InvokableLambda(func(ctx context.Context, input map[string]any) (output map[string]any, err error) {
		fmt.Println("====处理完成")
		return input, nil
	}))
	runnable, err := chain.Compile(ctx)
	if err != nil {
		fmt.Println("编译失败", err)
	}
	testCases := []map[string]any{
		{"language": "go", "task": "开发AI应用"},
		{"language": "go", "task": "开发AI应用"},
		{"language": "go", "task": "开发AI应用"},
	}
	for _, testCase := range testCases {
		content2, _ := runnable.Invoke(ctx, testCase)
		fmt.Println(content2)
	}

}
