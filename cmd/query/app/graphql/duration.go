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
	Name:        "Duration",
	Description: "时间区间参数，定义查询的有效时间区间",
	Fields: graphql.InputObjectConfigFieldMap{
		"start": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "开始时间，格式为2016-01-01 1310",
		},
		"end": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "结束时间，格式为2016-01-01 1310",
		},
		"step": &graphql.InputObjectFieldConfig{
			Type:         graphql.String,
			DefaultValue: "MINUTE",
			Description:  "分桶策略，默认是MINUTE，用于需要进行时间分桶的查询，如热力图，趋势图",
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
	var interval time.Duration
	if d.Step == "HOUR" {
		interval = time.Hour
	} else if d.Step == "DAY" {
		interval = time.Hour * 24
	} else {
		interval = time.Minute
	}
	return &spanstore.ThermoDynamicQueryParameters{
		BasicQueryParameters: spanstore.BasicQueryParameters{
			StartTimeMin: bq.StartTimeMin,
			StartTimeMax: bq.StartTimeMax,
		},
		TimeInterval:            interval,
		DurationInterval:        time.Millisecond * 100,
		DurationExtendBoundsMin: 0,
		DurationExtendBoundsMax: time.Millisecond * 3000,
	}, nil
}
