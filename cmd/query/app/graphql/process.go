package graphql

import (
	"github.com/graphql-go/graphql"
)

// type Process struct {
// 	ServiceName string     `json:"serviceName"`
// 	Tags        []KeyValue `json:"tags"`
// }

var GLProcessType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Process",
		Fields: graphql.Fields{
			"serviceName": &graphql.Field{
				Type: graphql.String,
			},
			"tags": &graphql.Field{
				Type: graphql.NewList(GLKeyValueType),
			},
		},
	},
)

var GLFlattenProcessType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "FlattenProcess",
		Fields: graphql.Fields{
			"key": &graphql.Field{
				Type: graphql.String,
			},
			"value": &graphql.Field{
				Type: GLProcessType,
			},
		},
	},
)
