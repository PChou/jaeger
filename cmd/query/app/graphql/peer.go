package graphql

import (
	"github.com/graphql-go/graphql"
)

type PeersInfo struct {
	Count int      `json:"count"`
	Peers []string `json:"peers"`
}

var GLPeersType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Peers",
		Fields: graphql.Fields{
			"count": &graphql.Field{
				Type: graphql.Int,
			},
			"peers": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
		},
	},
)
