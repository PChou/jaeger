package app

import (
//"github.com/graphql-go/graphql"
)

// var operationType = graphql.NewObject(
// 	graphql.ObjectConfig{
// 		Name: "Operation",
// 		Fields: graphql.Fields{
// 			"name": &graphql.Field{
// 				Type: graphql.String,
// 			},
// 		},
// 	},
// )

type Service struct {
	Name       string      `json:"name"`
	Operations []Operation `json:"operations"`
}

type Operation struct {
	Name string `json:"name"`
}

type GraphQLPostBody struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}
