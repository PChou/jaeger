package app

import (
	"time"

	"github.com/graphql-go/graphql"
	"github.com/jaegertracing/jaeger/storage/spanstore"
)

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

// Duration used to supply thermodynamic query
type Duration struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Step  string `json:"step"`
}

func (d Duration) ToThermoDynamicQueryParameters() (*spanstore.ThermoDynamicQueryParameters, error) {
	start, err := time.ParseInLocation("2006-1-2 15:4", d.Start, time.Now().Location())
	if err != nil {
		return nil, err
	}
	end, err := time.ParseInLocation("2006-1-2 15:4", d.End, time.Now().Location())
	if err != nil {
		return nil, err
	}
	//d.Step is hard code to "MINUTE"

	return &spanstore.ThermoDynamicQueryParameters{
		//ServiceName:
		//OperationName:    "callSample",
		StartTimeMin:            start,
		StartTimeMax:            end,
		TimeInterval:            time.Minute,
		DurationInterval:        time.Millisecond * 500,
		DurationExtendBoundsMin: 0,
		DurationExtendBoundsMax: time.Millisecond * 3000,
	}, nil
}

//Input Type
var durationType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "Duration",
	Description: "duration that passed in",
	Fields: graphql.InputObjectConfigFieldMap{
		"start": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"end": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"step": &graphql.InputObjectFieldConfig{
			Type:         graphql.String,
			DefaultValue: "MINUTE",
		},
	},
})

var nodeType = graphql.NewList(graphql.Int)
