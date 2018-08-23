package graphql

import (
	"github.com/graphql-go/graphql"
)

type Operation struct {
	Name string `json:"name"`
}

var GLOperationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Operation",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
