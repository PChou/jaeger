package graphql

import (
	"github.com/graphql-go/graphql"
)

type ServiceList struct {
	Services []string `json:"services"`
	Count    int      `json:"count"`
}

var GLServiceListType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "ServiceList",
		Description: "服务列表",
		Fields: graphql.Fields{
			"services": &graphql.Field{
				Type:        graphql.NewList(graphql.String),
				Description: "服务名称列表",
			},
			"count": &graphql.Field{
				Type:        graphql.Int,
				Description: "服务数量",
			},
		},
	},
)

// model/ext_reader.go contains the structed ServiceAvgResponseTime
var GLServiceAvgResponseTime = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "ServiceResponseTime",
		Description: "服务平均响应时间",
		Fields: graphql.Fields{
			"applicationName": &graphql.Field{
				Description: "服务所属的应用",
				Type:        graphql.String,
			},
			"serviceName": &graphql.Field{
				Description: "服务名称",
				Type:        graphql.String,
			},
			"value": &graphql.Field{
				Description: "响应时间(毫秒)",
				Type:        graphql.Float,
			},
		},
	},
)
