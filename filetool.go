package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cloudwego/eino/schema"
)

/*
优点：
流程完全透明，每一步 JSON 转换都看得见，新手好理解；
自由度极高：中间可以加自定义日志、特殊异常包装、参数校验、权限拦截等定制逻辑；
复杂工具、有特殊预处理 / 后处理逻辑时适合这种写法。
缺点：
重复模板代码多：每个工具都要写一遍 Unmarshal、Marshal；
容易漏处理错误（代码里 json.Marshal 全部忽略 err）；
新增工具就要重复拷贝一套 JSON 胶水代码，冗余。

你这段 FileReaderTool 是原生手动实现框架 Tool 接口，所有 JSON 序列化、反序列化、参数解析全自己手写；
前面 utils.NewTool 是框架提供的通用工具包装器，自动帮你抹平 JSON 转换胶水代码，只需要写纯业务函数

utils.NewTool 内部做了什么（底层封装逻辑）：
通过反射读取你的业务函数：
第一个入参结构体：TimeParams（入参模型）
返回第一个结构体：TimeResult（出参模型）
自动生成 Info() 方法：直接复用你传入的 schema.ToolInfo；
自动生成完整 InvokableRun 胶水逻辑（和你手动写的逻辑一模一样，只是封装在框架内部）：

优点：
业务代码干净，只关心业务逻辑，无重复 JSON 样板代码；
统一处理序列化错误，不用每次手动判断；
快速开发简单工具，一行生成工具实例。
缺点：
胶水逻辑被封装，看不到底层转换过程；
如果需要在参数解析前后加自定义拦截、特殊格式化、日志埋点，需要额外套一层 Lambda，不如手动实现灵活。
二、两者核心差异对照表
表格
对比项	手动实现 Tool（FileReaderTool）	utils.NewTool 包装函数
接口实现	自己手动写 Info + InvokableRun	框架内部自动实现接口
JSON 转换	手动 Unmarshal / Marshal	反射自动完成，无需手写
代码冗余	每个工具重复写 JSON 胶水代码	无重复模板代码
自定义能力	极高，可任意修改入参出参处理逻辑	基础场景够用，复杂预处理需要额外封装
适用场景	工具需要特殊校验、异常包装、权限控制、复杂转换	简单工具（获取时间、简单计算、查询），快速开发
错误处理	自己控制序列化失败逻辑	框架统一捕获序列化错误

简单轻量工具（获取时间、简单四则运算、查天气）→ 用 utils.NewTool，少写重复代码；
复杂定制工具（需要参数校验、路径拦截、自定义错误格式、打印详细入参日志、权限校验）→ 手动实现 Tool 接口，自由控制 InvokableRun 全流程；
学习理解框架原理：手写 FileReaderTool 更容易看懂工具和 LLM 交互的 JSON 流转逻辑
*/
type FileReaderTool struct {
}

type FileParams struct {
	FilePath string `json:"filepath"`
}
type FileResult struct {
	Content string `json:"content"`
	Error   error  `json:"error,omitempty"`
}

func (t *FileReaderTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "read_file",
		Desc: "读取文件内容",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"filepath": {
				Type:     schema.String,
				Required: true,
				Desc:     "文件路径"},
		}),
	}, nil
}
func (t *FileReaderTool) InvokableRun(ctx context.Context, argumengtJson string) (string, error) {
	var params FileParams
	if err := json.Unmarshal([]byte(argumengtJson), &params); err != nil {
		return "", err
	}
	filepath := params.FilePath
	content, err := os.ReadFile(filepath)

	if err != nil {
		result := &FileResult{Error: fmt.Errorf("文件读取失败 %v", err)}
		b, _ := json.Marshal(result)
		return string(b), err
	}

	result := &FileResult{Content: string(content)}
	contents, _ := json.Marshal(result)
	return string(contents), nil
}

func main() {
	fileTool := &FileReaderTool{}
	var params FileParams
	params.FilePath = "a.txt"
	b, _ := json.Marshal(params)
	fmt.Println(fileTool.InvokableRun(context.Background(), string(b)))
}
