package graphql

import (
	"github.com/graphql-go/graphql"
)

type ApplicationList struct {
	Applications []string `json:"applications"`
	Count        int      `json:"count"`
}

var GLApplicationListType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "ApplicationList",
		Description: "应用列表",
		Fields: graphql.Fields{
			"applications": &graphql.Field{
				Type:        graphql.NewList(graphql.String),
				Description: "应用名称，可作为应用的唯一标识",
			},
			"count": &graphql.Field{
				Type:        graphql.Int,
				Description: "应用数量",
			},
		},
	},
)

var GLApplicationThroughput = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ApplicationThroughput",
		Fields: graphql.Fields{
			"applicationName": &graphql.Field{
				Type:        graphql.String,
				Description: "应用名称",
			},
			"value": &graphql.Field{
				Type:        graphql.Float,
				Description: "应用在指定时间段内的每分钟请求次数",
			},
		},
	},
)
