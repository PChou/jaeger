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
var GLDependencyExt = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "DependencyExt",
		Description: "表示一个调用关系",
		Fields: graphql.Fields{
			"parent": &graphql.Field{
				Type:        graphql.String,
				Description: "调用者",
			},
			"parentLayer": &graphql.Field{
				Type:        graphql.String,
				Description: "调用者的类型，比如HTTP，DB，CACHE",
			},
			"parentComponent": &graphql.Field{
				Type:        graphql.String,
				Description: "调用者组件名称，比如Tomcat",
			},
			"child": &graphql.Field{
				Type:        graphql.String,
				Description: "被调用者",
			},
			"childLayer": &graphql.Field{
				Type:        graphql.String,
				Description: "被调用者的类型，比如HTTP，DB，CACHE",
			},
			"childComponent": &graphql.Field{
				Type:        graphql.String,
				Description: "被调用者组件名称，比如Tomcat",
			},
			"callCount": &graphql.Field{
				Type:        graphql.Int,
				Description: "累计调用次数",
			},
		},
	},
)
