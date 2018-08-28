package graphql

import (
	"github.com/graphql-go/graphql"
)

var GLNodeType = graphql.NewList(graphql.Int)

var GLThermodynamicType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "ThermoDynamic",
		Description: "响应时间热力图",
		Fields: graphql.Fields{
			"responseTimeStep": &graphql.Field{
				Description: "响应时间的区间，单位是ms",
				Type:        graphql.Int,
			},
			"nodes": &graphql.Field{
				Description: "点阵图，这是一个整型二维数组",
				Type:        graphql.NewList(GLNodeType),
			},
		},
	},
)

// model/ext_reader.go contains the structed ThermodynamicType
