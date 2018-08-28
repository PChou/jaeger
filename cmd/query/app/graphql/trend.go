package graphql

import (
	"github.com/graphql-go/graphql"
)

type Trends struct {
	TrendList []int `json:"trendList"`
}

var GLTrendListType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "TrendList",
		Description: "趋势图",
		Fields: graphql.Fields{
			"trendList": &graphql.Field{
				Type:        graphql.NewList(graphql.Int),
				Description: "趋势图",
			},
		},
	},
)
