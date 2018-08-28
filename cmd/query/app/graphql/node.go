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
		Name: "NodeAvgThroughput",
		Fields: graphql.Fields{
			"node": &graphql.Field{
				Type: graphql.String,
			},
			"value": &graphql.Field{
				Type: graphql.Float,
			},
		},
	},
)
