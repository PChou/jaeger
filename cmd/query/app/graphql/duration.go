package graphql

import (
	"time"

	"github.com/graphql-go/graphql"
	"github.com/jaegertracing/jaeger/storage/spanstore"
)

// Duration used to supply query associated with time
type Duration struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Step  string `json:"step"`
}

//Input Type
var GLDurationType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "Duration",
	//Description: "duration that passed in",
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

func (d Duration) ToThermoDynamicQueryParameters() (*spanstore.ThermoDynamicQueryParameters, error) {
	start, err := ParseSkyWalkingTimeFormat(d.Start)
	if err != nil {
		return nil, err
	}
	end, err := ParseSkyWalkingTimeFormat(d.End)
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

func (d Duration) ToApplicationQueryParameters() (*spanstore.ApplicationQueryParameter, error) {
	start, err := ParseSkyWalkingTimeFormat(d.Start)
	if err != nil {
		return nil, err
	}
	end, err := ParseSkyWalkingTimeFormat(d.End)
	if err != nil {
		return nil, err
	}

	return &spanstore.ApplicationQueryParameter{
		StartTimeMin: start,
		StartTimeMax: end,
	}, nil
}

func (d Duration) ToTrendsQueryParameters() (*spanstore.TrendsQueryParameters, error) {
	start, err := ParseSkyWalkingTimeFormat(d.Start)
	if err != nil {
		return nil, err
	}
	end, err := ParseSkyWalkingTimeFormat(d.End)
	if err != nil {
		return nil, err
	}

	return &spanstore.TrendsQueryParameters{
		StartTimeMin: start,
		StartTimeMax: end,
	}, nil
}
