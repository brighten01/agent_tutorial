package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/schema"
)

type CalculatorTool struct {
}

type CalculatorParams struct {
	Operation string  `json:"operation"`
	A         float64 `json:"a"`
	B         float64 `json:"b"`
}

type CalculatorResult struct {
	Result float64 `json:"result"`
	Error  string  `json:"error,omitempty"`
}

// 用户输入->AI 理解意图->选择 Tool->执行Tool->获取结果->AI生成回复
//
//	查天气-> "weather_tool"->"晴天 25度"

func (t *CalculatorTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	// 返回工具元信息，供 AI 感知该工具的名称、描述和参数定义
	return &schema.ToolInfo{
		Name: "calculator",
		Desc: "执行基本数学运算（加、减、乘、除）",
		// 定义工具接受的输入参数列表
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			// operation 指定运算类型，非必填
			"operation": {
				Type: "string",
				Desc: "运算类型:add(加), substract(减) ,multiply(乘), divide(除)",
			},
			// a 为第一个操作数，必填
			"a": {
				Type:     "number",
				Desc:     "第一个数字",
				Required: true,
			},
			// b 为第二个操作数，必填
			"b": {
				Type:     "number",
				Desc:     "第二个数字",
				Required: true,
			},
		}),
	}, nil
}

func (t *CalculatorTool) InvokableRun(ctx context.Context, argumentsJson string) (string, error) {
	var params CalculatorParams
	if err := json.Unmarshal([]byte(argumentsJson), &params); err != nil {
		return "", fmt.Errorf("解析参数失败 ,%v", err)
	}
	var result float64
	switch params.Operation {
	case "add":
		result = params.A + params.B
	case "substract":
		result = params.A - params.B
	case "multiply":
		result = params.A * params.B
	case "divide":
		if params.B == 0 {
			resultJson, _ := json.Marshal(CalculatorResult{
				Error: "除数不能为0",
			})
			return string(resultJson), fmt.Errorf("除数不能为0")
		}
		result = params.A / params.B
	default:
		resultJson, _ := json.Marshal(CalculatorResult{
			Error: "不支持的运算类型",
		})
		return string(resultJson), fmt.Errorf("除数不能为0")
	}

	resultJson, err := json.Marshal(CalculatorResult{
		Result: result,
	})

	if err != nil {
		return "", err
	}
	return string(resultJson), nil
}

func main() {
	ctx := context.Background()
	calculator := &CalculatorTool{}
	testCases := []struct {
		operation string
		a, b      float64
	}{
		{"add", 1, 2},
		{"substract", 3, 4},
		{"multiply", 5, 6},
		{"divide", 7, 8},
		{"divide", 7, 0},
	}
	for _, testCase := range testCases {
		paramsJson, _ := json.Marshal(CalculatorParams{
			Operation: testCase.operation,
			A:         testCase.a,
			B:         testCase.b,
		})

		result, err := calculator.InvokableRun(ctx, string(paramsJson))
		if err != nil {
			fmt.Printf("执行失败: %v\n", err)
			continue
		}

		fmt.Println(result)
	}
}
