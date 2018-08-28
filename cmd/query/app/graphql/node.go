package graphql

import (
	"github.com/graphql-go/graphql"
)

type Node struct {
	ApplicationName string `json:"applicatoinName"`
	Name            string `json:"name"`
	//OS               string `json:"os"`
}

var GLNodeAvgThroughput = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "NodeAvgThroughput",
		Description: "节点（服务器）在指定时间区间内平均每分钟请求数",
		Fields: graphql.Fields{
			"node": &graphql.Field{
				Description: "节点（服务器）名称，可作为节点（服务器）的Id使用",
				Type:        graphql.String,
			},
			"value": &graphql.Field{
				Description: "平均每分钟请求数",
				Type:        graphql.Float,
			},
		},
	},
)
