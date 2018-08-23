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
