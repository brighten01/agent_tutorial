package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	ctx := context.Background()
	apiKey := "ark-94d2cd4b-30a4-44d9-99a7-afa4469d9009-09252"
	endpointID := "ep-20260715134349-f7jrm"
	url := "https://ark.cn-beijing.volces.com/api/v3/embeddings/multimodal"

	reqBody := map[string]any{
		"model": endpointID,
		"input": []any{
			map[string]any{
				"type": "text",
				"text": "测试文本",
			},
			map[string]any{
				"type": "text",
				"text": "测试文本11",
			},
		},
	}

	bodyData, _ := json.MarshalIndent(reqBody, "", "  ")
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("请求失败：", err)
		return
	}
	defer resp.Body.Close()

	resBody, _ := io.ReadAll(resp.Body)
	fmt.Println("调用结果：\n", string(resBody))
}
