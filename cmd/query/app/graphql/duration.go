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

func (d Duration) ToBasicQueryParameters() (*spanstore.BasicQueryParameters, error) {
	start, err := ParseSkyWalkingTimeFormat(d.Start)
	if err != nil {
		return nil, err
	}
	end, err := ParseSkyWalkingTimeFormat(d.End)
	if err != nil {
		return nil, err
	}

	return &spanstore.BasicQueryParameters{
		StartTimeMin: start,
		StartTimeMax: end,
	}, nil
}

func (d Duration) ToThermoDynamicQueryParameters() (*spanstore.ThermoDynamicQueryParameters, error) {
	bq, err := d.ToBasicQueryParameters()
	if err != nil {
		return nil, err
	}
	//TODO: d.Step is hard code to "MINUTE"
	return &spanstore.ThermoDynamicQueryParameters{
		BasicQueryParameters: spanstore.BasicQueryParameters{
			StartTimeMin: bq.StartTimeMin,
			StartTimeMax: bq.StartTimeMax,
		},
		TimeInterval:            time.Minute,
		DurationInterval:        time.Millisecond * 100,
		DurationExtendBoundsMin: 0,
		DurationExtendBoundsMax: time.Millisecond * 3000,
	}, nil
}
