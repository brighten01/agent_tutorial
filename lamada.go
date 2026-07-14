package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cloudwego/eino/compose"
)

// 1、链式调用函数
// 2、创建chain 定义输出 文本->文本
// 3、AppendLambda() 函数定义匿名函数输入 返回
// 4、编译 Compile
// 5、编译通过执行Invoke
// 6、输出结果

func main() {
	ctx := context.Background()
	chain := compose.NewChain[string, string]()
	chain.
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
			fmt.Printf("步骤1: 输入=%s\\n", input)
			result := strings.ToUpper(input)
			fmt.Printf("步骤1:输出 = %s \\n", result)
			return result, nil
		})).
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
			fmt.Printf("步骤2: 输入=%s\\n", input)
			result := "处理结果1111" + input
			fmt.Printf("步骤2:输出 = %s \\n", result)
			return result, nil
		}))

	runnable, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("编译失败: %v", err)
	}

	output, _ := runnable.Invoke(ctx, "hello")
	fmt.Println(output)
}
