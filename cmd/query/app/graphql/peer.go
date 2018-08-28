package graphql

import (
	"github.com/graphql-go/graphql"
)

type PeerList struct {
	Count int      `json:"count"`
	Peers []string `json:"peers"`
}

var GLPeersType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "PeerList",
		Description: "对端的服务地址列表，用于表达数据库和缓存这类服务",
		Fields: graphql.Fields{
			"count": &graphql.Field{
				Type:        graphql.Int,
				Description: "地址数量",
			},
			"peers": &graphql.Field{
				Type:        graphql.NewList(graphql.String),
				Description: "对端的服务地址列表，是ip和端口的组合",
			},
		},
	},
)
