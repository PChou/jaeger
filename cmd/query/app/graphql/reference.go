package graphql

import (
	"github.com/graphql-go/graphql"
)

// json/model.Reference
// type Reference struct {
// 	RefType ReferenceType `json:"refType"`
// 	TraceID TraceID       `json:"traceID"`
// 	SpanID  SpanID        `json:"spanID"`
// }
var GLReferenceType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Reference",
		Fields: graphql.Fields{
			"refType": &graphql.Field{
				Type: graphql.String,
			},
			"traceID": &graphql.Field{
				Type: graphql.String,
			},
			"spanID": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
