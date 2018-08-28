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
		Name: "ServiceList",
		Fields: graphql.Fields{
			"services": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"count": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

// model/ext_reader.go contains the structed ServiceAvgResponseTime
var GLServiceAvgResponseTime = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ServiceResponseTime",
		Fields: graphql.Fields{
			"applicationName": &graphql.Field{
				Type: graphql.String,
			},
			"serviceName": &graphql.Field{
				Type: graphql.String,
			},
			"value": &graphql.Field{
				Type: graphql.Float,
			},
		},
	},
)
