package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

type Restaurant struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Location string   `json:"location"`
	Cuisine  string   `json:"cuisine"`
	Rating   float64  `json:"rating"`
	Tags     []string `json:"tags"`
}

type Dish struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
	Spicy bool   `json:"spicy"`
	Desc  string `json:"desc"`
}

var datas = []Restaurant{
	{
		ID:       "r1",
		Name:     "川湘阁",
		Location: "北京",
		Cuisine:  "川菜",
		Rating:   4.8,
		Tags:     []string{"辣", "正宗", "环境好"},
	},
	{
		ID:       "r2",
		Name:     "眉州东坡酒楼",
		Location: "北京",
		Cuisine:  "粤菜",
		Rating:   4.5,
		Tags:     []string{"清淡", "海鲜", "精致"},
	},
	{
		ID:       "r3",
		Name:     "麻辣香锅",
		Location: "北京",
		Cuisine:  "粤菜",
		Rating:   4.5,
		Tags:     []string{"辣", "实惠", "分量足"},
	},
}

var dishes = map[string][]Dish{
	"r1": {
		{
			Name:  "水煮鱼",
			Price: 88,
			Spicy: true,
			Desc:  "鲜嫩鱼肉，麻辣鲜香",
		},
		{
			Name:  "宫保鸡丁",
			Price: 15,
			Spicy: true,
			Desc:  "经典川菜，香辣可口",
		},
		{
			Name:  "蒜泥白肉",
			Price: 39,
			Spicy: false,
			Desc:  "肥而不腻，蒜香浓郁",
		},
	},
	"r2": {
		{
			Name:  "清蒸鲈鱼",
			Price: 30,
			Spicy: false,
			Desc:  "鲜嫩多汁，原汁原味",
		},
		{
			Name:  "白切鸡",
			Price: 30,
			Spicy: false,
			Desc:  "皮爽肉滑，鸡味浓郁",
		},
		{
			Name:  "广式烧鹅",
			Price: 30,
			Spicy: false,
			Desc:  "皮脆肉嫩，香味四溢",
		},
	},

	"r3": {
		{
			Name:  "麻辣香锅",
			Price: 30,
			Spicy: true,
			Desc:  "自选食材，麻辣鲜香",
		},
		{
			Name:  "干锅牛蛙",
			Price: 30,
			Spicy: true,
			Desc:  "肉质细嫩，香辣入味",
		},
		{
			Name:  "毛血旺",
			Price: 58,
			Spicy: true,
			Desc:  "麻辣鲜香，食材丰富",
		},
	},
}

func main() {
	ctx := context.Background()
	chatmodel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  os.Getenv("DEEPSEEK_API_KEY"),
		Model:   "deepseek-chat",
		BaseURL: "https://api.deepseek.com",
	})
	if err != nil {
		log.Fatalf("failed to create chat model: %v", err)
	}
	restaurantTool := utils.NewTool(&schema.ToolInfo{
		Name: "query_restaurant",
		Desc: "餐厅查询工具",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"location": {
				Type:     "string",
				Desc:     "城市位置例如北京、上海",
				Required: false,
			},
			"cuisine": {
				Type:     "string",
				Desc:     "菜系 ，例如 川菜、粤菜、湘菜",
				Required: false,
			},

			//3 failed to generate content: [NodeRunError] failed to create chat completion: HTTP 400: Bad request
			//{"error":{"message":"Invalid schema for function 'query_restaurant': \"bool\" is not valid under any of the schemas listed in the 'anyOf' keyword","type":"invalid_request_error","param":null,"code":"invalid_request_error"}}
			//-------------------- type 一定要写正确 boolean
			//const (
			//	Object  DataType = "object"
			//	Number  DataType = "number"
			//	Integer DataType = "integer"
			//	String  DataType = "string"
			//	Array   DataType = "array"
			//	Null    DataType = "null"
			//	Boolean DataType = "boolean" //和go 类型不同
			//)

			"spicy": {
				Type: "boolean", //
				Desc: "是否要辣的",
			},
		}),
	}, func(ctx context.Context, params map[string]any) (string, error) {
		fmt.Printf("[工具执行] query_restaurant\\n")
		var result []Restaurant
		location, _ := params["location"].(string)
		cuisine, _ := params["cuisine"].(string)
		spicy, _ := params["spicy"].(bool)

		for _, restaurant := range datas {
			match := true
			if location != "" && restaurant.Location != location {
				match = false
			}
			if cuisine != "" && restaurant.Cuisine != cuisine {
				match = false
			}
			if spicy {
				hasSpicy := false
				for _, tag := range restaurant.Tags {
					if tag == "辣" {
						hasSpicy = true
						break
					}
				}

				if !hasSpicy {
					match = false
				}
			}

			if match {
				result = append(result, restaurant)
			}
		}
		resultJson, _ := json.Marshal(result)
		return string(resultJson), nil
	})

	//菜品
	dishTool := utils.NewTool(&schema.ToolInfo{
		Name: "query_dishes",
		Desc: "查询指定餐厅（需要餐厅ID",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"restaurant_id": {
				Type:     "string",
				Desc:     "餐厅ID",
				Required: true,
			},
		}),
	}, func(ctx context.Context, params map[string]any) (string, error) {
		fmt.Printf("[工具执行] query_dishes\\n")
		restaurantID, _ := params["restaurant_id"].(string)
		if dishes, ok := dishes[restaurantID]; ok {
			resultJson, _ := json.Marshal(dishes)
			return string(resultJson), nil
		}
		return "", fmt.Errorf("restaurant not found")
	})

	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatmodel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{
				restaurantTool,
				dishTool,
			},
		},
	})

	messages := []*schema.Message{
		schema.UserMessage("我在北京想吃辣一点的，给我推荐几家餐厅和特色菜"),
	}
	fmt.Println("=====用户在北京，想吃辣一点的，给我推荐几家北京的特色菜=====")

	content, err := agent.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("failed to generate content: %v", err)
	}
	fmt.Println("===Agent 回答======")

	fmt.Println(content.Content)

}
