package graphql

import (
	"github.com/graphql-go/graphql"
)

var GLNodeType = graphql.NewList(graphql.Int)

var GLThermodynamicType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ThermoDynamic",
		Fields: graphql.Fields{
			"responseTimeStep": &graphql.Field{
				Type: graphql.Int,
			},
			"nodes": &graphql.Field{
				Type: graphql.NewList(GLNodeType),
			},
		},
	},
)

// model/ext_reader.go contains the structed ThermodynamicType
