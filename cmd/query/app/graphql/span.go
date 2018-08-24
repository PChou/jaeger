package graphql

import (
	"github.com/graphql-go/graphql"
)

var GLKeyValueType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "KeyValue",
		Fields: graphql.Fields{
			"key": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
			"value": &graphql.Field{
				Type: graphql.NewScalar(graphql.ScalarConfig{
					Name: "TagValue",
					Serialize: func(value interface{}) interface{} {
						return value
					},
				}),
			},
		},
	},
)

var GLLogType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Log",
		Fields: graphql.Fields{
			"timestamp": &graphql.Field{
				Type: graphql.Int,
			},
			"fields": &graphql.Field{
				Type: graphql.NewList(GLKeyValueType),
			},
		},
	},
)

// json/model.Trace(ui.Trace)
// type Span struct {
// 	TraceID       TraceID     `json:"traceID"`
// 	SpanID        SpanID      `json:"spanID"`
// 	ParentSpanID  SpanID      `json:"parentSpanID,omitempty"` // deprecated
// 	Flags         uint32      `json:"flags,omitempty"`
// 	OperationName string      `json:"operationName"`
// 	References    []Reference `json:"references"`
// 	StartTime     uint64      `json:"startTime"` // microseconds since Unix epoch
// 	Duration      uint64      `json:"duration"`  // microseconds
// 	Tags          []KeyValue  `json:"tags"`
// 	Logs          []Log       `json:"logs"`
// 	ProcessID     ProcessID   `json:"processID,omitempty"`
// 	Process       *Process    `json:"process,omitempty"`
// 	Warnings      []string    `json:"warnings"`
// }

var GLSpanType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Span",
		Fields: graphql.Fields{
			"traceID": &graphql.Field{
				Type: graphql.String,
			},
			"spanID": &graphql.Field{
				Type: graphql.String,
			},
			"parentSpanID": &graphql.Field{
				Type: graphql.String,
			},
			"flags": &graphql.Field{
				Type: graphql.Int,
			},
			"operationName": &graphql.Field{
				Type: graphql.String,
			},
			"references": &graphql.Field{
				Type: graphql.NewList(GLReferenceType),
			},
			"startTime": &graphql.Field{
				Type: graphql.Float,
			},
			"duration": &graphql.Field{
				Type: graphql.Int,
			},
			"tags": &graphql.Field{
				Type: graphql.NewList(GLKeyValueType),
			},
			"logs": &graphql.Field{
				Type: graphql.NewList(GLLogType),
			},
			"processID": &graphql.Field{
				Type: graphql.String,
			},
			"process": &graphql.Field{
				Type: GLProcessType,
			},
			"warnings": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
		},
	},
)
