package graphql

import (
	"time"

	"github.com/graphql-go/graphql"
	ui "github.com/jaegertracing/jaeger/model/json"
	"github.com/jaegertracing/jaeger/storage/spanstore"
)

var GLTraceType2 = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Trace",
		Fields: graphql.Fields{
			"key": &graphql.Field{
				Type: graphql.String,
			},
			"isError": &graphql.Field{
				Type: graphql.Boolean,
			},
			"duration": &graphql.Field{
				Type: graphql.Int,
			},
			"start": &graphql.Field{
				Type: graphql.String,
			},
			"traceIds": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"operationNames": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
		},
	},
)

// json/model.Trace(ui.Trace)
// type Trace struct {
// 	TraceID   TraceID               `json:"traceID"`
// 	Spans     []Span                `json:"spans"`
// 	Processes map[ProcessID]Process `json:"processes"`
// 	Warnings  []string              `json:"warnings"`
// }
var GLTraceType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Trace",
		Fields: graphql.Fields{
			"traceID": &graphql.Field{
				Type: graphql.String,
			},
			"spans": &graphql.Field{
				Type: graphql.NewList(GLSpanType),
			},
			//TODO don't know how to export Processes which is a map
			//here workaround use array
			"flattenProcesses": &graphql.Field{
				Type: graphql.NewList(GLFlattenProcessType),
			},
			"warnings": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
		},
	},
)

type TraceList struct {
	Total  int         `json:"total"`
	Traces []*ui.Trace `json:"traces"`
}

var GLTraceListType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "TraceList",
		Fields: graphql.Fields{
			"total": &graphql.Field{
				Type: graphql.Int,
			},
			"traces": &graphql.Field{
				Type: graphql.NewList(GLTraceType),
			},
		},
	},
)

type TraceQueryCondition struct {
	ApplicationId    string   `json:"applicationId"`
	OperationName    string   `json:"operationName"`
	MaxTraceDuration int      `json:"maxTraceDuration"` //ms
	MinTraceDuration int      `json:"minTraceDuration"` //ms
	QueryOrder       string   `json:"queryOrder"`
	TraceState       string   `json:"traceState"`
	QueryDuration    Duration `json:"queryDuration"`
	Paging           struct {
		NeedTotal bool `json:"needTotal"`
		PageNum   int  `json:"pageNum"`
		PageSize  int  `json:"pageSize"`
	} `json:"paging"`
}

func (c TraceQueryCondition) ToTraceQueryParameters() (*spanstore.TraceQueryParameters, error) {
	start, err := ParseSkyWalkingTimeFormat(c.QueryDuration.Start)
	if err != nil {
		return nil, err
	}
	end, err := ParseSkyWalkingTimeFormat(c.QueryDuration.End)
	if err != nil {
		return nil, err
	}

	return &spanstore.TraceQueryParameters{
		ServiceName:   c.ApplicationId,
		OperationName: c.OperationName,
		StartTimeMin:  start,
		StartTimeMax:  end,
		DurationMin:   time.Duration(c.MinTraceDuration) * time.Millisecond,
		DurationMax:   time.Duration(c.MaxTraceDuration) * time.Millisecond,
		NumTraces:     20, //at most get 5000
	}, nil
}

var GLTraceQueryConditionType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "TraceQueryCondition",
	//Description: "",
	Fields: graphql.InputObjectConfigFieldMap{
		"applicationId": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"operationName": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"maxTraceDuration": &graphql.InputObjectFieldConfig{
			Type: graphql.Int,
		},
		"minTraceDuration": &graphql.InputObjectFieldConfig{
			Type: graphql.Int,
		},
		"queryOrder": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"traceState": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"paging": &graphql.InputObjectFieldConfig{
			Type: graphql.NewInputObject(graphql.InputObjectConfig{
				Name: "Paging",
				Fields: graphql.InputObjectConfigFieldMap{
					"needTotal": &graphql.InputObjectFieldConfig{
						Type: graphql.Boolean,
					},
					"pageNum": &graphql.InputObjectFieldConfig{
						Type: graphql.Int,
					},
					"pageSize": &graphql.InputObjectFieldConfig{
						Type: graphql.Int,
					},
				},
			}),
		},
		"queryDuration": &graphql.InputObjectFieldConfig{
			Type: GLDurationType,
		},
	},
})
