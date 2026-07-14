package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cloudwego/eino/schema"
)

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
