package graphql

import (
	"github.com/graphql-go/graphql"
)

type ApplicationList struct {
	Applications []string `json:"applications"`
	Count        int      `json:"count"`
}

// var GLApplicationType = graphql.NewObject(
// 	graphql.ObjectConfig{
// 		Name: "Application",
// 		Fields: graphql.Fields{
// 			"name": &graphql.Field{
// 				Type: graphql.String,
// 			},
// 			// "operations": &graphql.Field{
// 			// 	Type: graphql.NewList(gl.GLOperationType),
// 			// 	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 			// 		if service, ok := p.Source.(Service); ok {
// 			// 			operations, err := aH.spanReader.GetOperations(service.Name)
// 			// 			if err != nil {
// 			// 				return nil, err
// 			// 			}
// 			// 			operationWrap := make([]gl.Operation, 0)
// 			// 			for _, operation := range operations {
// 			// 				operationWrap = append(operationWrap, gl.Operation{Name: operation})
// 			// 			}
// 			// 			return operationWrap, nil
// 			// 		}
// 			// 		return []interface{}{}, nil
// 			// 	},
// 			// },
// 		},
// 	},
// )

var GLApplicationListType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ApplicationList",
		Fields: graphql.Fields{
			"applications": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"count": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

var GLApplicationThroughput = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ApplicationThroughput",
		Fields: graphql.Fields{
			"applicationName": &graphql.Field{
				Type: graphql.String,
			},
			"value": &graphql.Field{
				Type: graphql.Float,
			},
		},
	},
)
