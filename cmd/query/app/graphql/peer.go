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
		Name: "PeerList",
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
