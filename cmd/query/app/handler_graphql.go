package app

import (
	"time"

	"github.com/graphql-go/graphql"
	gl "github.com/jaegertracing/jaeger/cmd/query/app/graphql"
	"github.com/jaegertracing/jaeger/model"
	ui "github.com/jaegertracing/jaeger/model/json"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func makeApplicationList(handler *APIHandler) *graphql.Field {
	return &graphql.Field{
		Type: gl.GLApplicationListType,
		Args: graphql.FieldConfigArgument{
			"duration": &graphql.ArgumentConfig{
				Type: gl.GLDurationType,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var durationParams gl.Duration
			err := mapstructure.Decode(p.Args["duration"], &durationParams)
			if err != nil {
				return nil, err
			}
			extReader := handler.spanReader.(spanstore.ExtReader)
			params, err := durationParams.ToBasicQueryParameters()
			if err != nil {
				return nil, err
			}
			applications, err := extReader.GetApplications(params)
			if err != nil {
				return []interface{}{}, err
			}

			return gl.ApplicationList{Applications: applications, Count: len(applications)}, nil
		},
	}
}

func makeServiceList(handler *APIHandler) *graphql.Field {
	return &graphql.Field{
		Type: gl.GLServiceListType,
		Args: graphql.FieldConfigArgument{
			"duration": &graphql.ArgumentConfig{
				Type: gl.GLDurationType,
			},
			"applicationName": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var durationParams gl.Duration
			err := mapstructure.Decode(p.Args["duration"], &durationParams)
			if err != nil {
				return nil, err
			}
			applicationName, ok := p.Args["applicationName"].(string)
			if !ok {
				return nil, errors.New("applicationName is not a valid String")
			}
			extReader := handler.spanReader.(spanstore.ExtReader)
			bqp, err := durationParams.ToBasicQueryParameters()
			if err != nil {
				return nil, err
			}
			services, err := extReader.GetLayerServices(&spanstore.LayerTypeQueryParameters{
				BasicQueryParameters: spanstore.BasicQueryParameters{
					StartTimeMin: bqp.StartTimeMin,
					StartTimeMax: bqp.StartTimeMax,
				},
				ApplicationName: applicationName,
				Layer:           "HTTP",
				Type:            "Entry",
				By:              "operationName",
			})
			if err != nil {
				return nil, err
			}
			return gl.ServiceList{Services: services, Count: len(services)}, nil
		},
	}
}

func makePeerList(handler *APIHandler, layer string) *graphql.Field {
	return &graphql.Field{
		Type: gl.GLPeersType,
		Args: graphql.FieldConfigArgument{
			"duration": &graphql.ArgumentConfig{
				Type: gl.GLDurationType,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var durationParams gl.Duration
			err := mapstructure.Decode(p.Args["duration"], &durationParams)
			if err != nil {
				return nil, err
			}
			extReader := handler.spanReader.(spanstore.ExtReader)
			bqp, err := durationParams.ToBasicQueryParameters()
			if err != nil {
				return nil, err
			}
			peers, err := extReader.GetLayerServices(&spanstore.LayerTypeQueryParameters{
				BasicQueryParameters: spanstore.BasicQueryParameters{
					StartTimeMin: bqp.StartTimeMin,
					StartTimeMax: bqp.StartTimeMax,
				},
				Layer: layer,
				Type:  "Exit",
				By:    "flattenTags.peer",
			})
			if err != nil {
				return nil, err
			}
			return gl.PeerList{Peers: peers, Count: len(peers)}, nil
		},
	}
}

func makeThermodynamic(handler *APIHandler) *graphql.Field {
	return &graphql.Field{
		Type: gl.GLThermodynamicType,
		Args: graphql.FieldConfigArgument{
			"duration": &graphql.ArgumentConfig{
				//DefaultValue:
				//Description:
				Type: gl.GLDurationType,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var durationParams gl.Duration
			err := mapstructure.Decode(p.Args["duration"], &durationParams)
			if err != nil {
				return nil, err
			}
			extReader := handler.spanReader.(spanstore.ExtReader)
			params, err := durationParams.ToThermoDynamicQueryParameters()
			if err != nil {
				return nil, err
			}
			td, err := extReader.GetThermoDynamic(params)
			if err != nil {
				return nil, err
			}
			return td, err
		},
	}
}

func makeTopSlowService(handler *APIHandler) *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(gl.GLServiceAvgResponseTime),
		Args: graphql.FieldConfigArgument{
			"duration": &graphql.ArgumentConfig{
				//DefaultValue:
				//Description:
				Type: gl.GLDurationType,
			},
			"applicationName": &graphql.ArgumentConfig{
				Type:         graphql.String,
				DefaultValue: "",
			},
			"topN": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				DefaultValue: 10,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var durationParams gl.Duration
			err := mapstructure.Decode(p.Args["duration"], &durationParams)
			if err != nil {
				return nil, err
			}
			topN, ok := p.Args["topN"].(int)
			if !ok {
				return nil, errors.New("topN is not a valid Int")
			}
			applicationName, ok := p.Args["applicationName"].(string)
			if !ok {
				return nil, errors.New("applicationName is not a valid Int")
			}
			extReader := handler.spanReader.(spanstore.ExtReader)
			bqp, err := durationParams.ToBasicQueryParameters()
			if err != nil {
				return nil, err
			}
			td, err := extReader.GetServiceTopResponseTime(&spanstore.ServiceTopResponseTimeQueryParameters{
				BasicQueryParameters: spanstore.BasicQueryParameters{
					StartTimeMin: bqp.StartTimeMin,
					StartTimeMax: bqp.StartTimeMax,
				},
				ApplicationName: applicationName,
				Top:             topN,
			})
			if err != nil {
				return nil, err
			}
			return td, err
		},
	}
}

func makeApplicationTopThroughput(handler *APIHandler) *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(gl.GLApplicationThroughput),
		Args: graphql.FieldConfigArgument{
			"duration": &graphql.ArgumentConfig{
				Type: gl.GLDurationType,
			},
			"topN": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				DefaultValue: 10,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var durationParams gl.Duration
			err := mapstructure.Decode(p.Args["duration"], &durationParams)
			if err != nil {
				return nil, err
			}
			topN, ok := p.Args["topN"].(int)
			if !ok {
				return nil, errors.New("topN is not a valid Int")
			}

			extReader := handler.spanReader.(spanstore.ExtReader)
			bqp, err := durationParams.ToBasicQueryParameters()
			if err != nil {
				return nil, err
			}
			td, err := extReader.GetApplicationTopThroughput(&spanstore.TopThroughputQueryParameters{
				BasicQueryParameters: spanstore.BasicQueryParameters{
					StartTimeMin: bqp.StartTimeMin,
					StartTimeMax: bqp.StartTimeMax,
				},
				Top: topN,
			})
			if err != nil {
				return nil, err
			}
			return td, err
		},
	}
}

func makeServerTopThroughput(handler *APIHandler) *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(gl.GLNodeAvgThroughput),
		Args: graphql.FieldConfigArgument{
			"duration": &graphql.ArgumentConfig{
				Type: gl.GLDurationType,
			},
			"applicationName": &graphql.ArgumentConfig{
				Type:         graphql.String,
				DefaultValue: "",
			},
			"topN": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				DefaultValue: 10,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var durationParams gl.Duration
			err := mapstructure.Decode(p.Args["duration"], &durationParams)
			if err != nil {
				return nil, err
			}
			//topN, ok := p.Args["topN"].(int)
			// if !ok {
			// 	return nil, errors.New("topN is not a valid Int")
			// }
			applicationName, ok := p.Args["applicationName"].(string)
			if !ok {
				return nil, errors.New("applicationName is not a valid Int")
			}
			extReader := handler.spanReader.(spanstore.ExtReader)
			bqp, err := durationParams.ToBasicQueryParameters()
			if err != nil {
				return nil, err
			}
			td, err := extReader.GetNodeTopThroughput(&spanstore.TopThroughputQueryParameters{
				BasicQueryParameters: spanstore.BasicQueryParameters{
					StartTimeMin: bqp.StartTimeMin,
					StartTimeMax: bqp.StartTimeMax,
				},
				ApplicationName: applicationName,
			})
			if err != nil {
				return nil, err
			}
			return td, err
		},
	}
}

func makeGLNode(handler *APIHandler) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Node",
			Fields: graphql.Fields{
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"os": &graphql.Field{
					Type: graphql.String,
				},
				"throughputTrends": &graphql.Field{
					Type: graphql.NewList(graphql.Int),
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						var durationParams gl.Duration
						err := mapstructure.Decode(p.Info.VariableValues["duration"], &durationParams)
						if err != nil {
							return nil, err
						}
						applicationName, ok := p.Info.VariableValues["applicationName"].(string)
						if !ok {
							return nil, errors.New("applicationName is not a valid String")
						}
						bqp, err := durationParams.ToBasicQueryParameters()
						if err != nil {
							return nil, err
						}
						if node, ok := p.Source.(gl.Node); ok {
							extReader := handler.spanReader.(spanstore.ExtReader)
							return extReader.GetThroughputTrends(&spanstore.ThroughputQueryParameters{
								BasicQueryParameters: spanstore.BasicQueryParameters{
									StartTimeMin: bqp.StartTimeMin,
									StartTimeMax: bqp.StartTimeMax,
								},
								ApplicationName: applicationName,
								Instance:        node.Name,
								TimeInterval:    time.Minute,
							})
						}
						return nil, nil
					},
				},
				"responseTimeTrends": &graphql.Field{
					Type: graphql.NewList(graphql.Float),
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						var durationParams gl.Duration
						err := mapstructure.Decode(p.Info.VariableValues["duration"], &durationParams)
						if err != nil {
							return nil, err
						}
						applicationName, ok := p.Info.VariableValues["applicationName"].(string)
						if !ok {
							return nil, errors.New("applicationName is not a valid String")
						}
						bqp, err := durationParams.ToBasicQueryParameters()
						if err != nil {
							return nil, err
						}
						if node, ok := p.Source.(gl.Node); ok {
							extReader := handler.spanReader.(spanstore.ExtReader)
							return extReader.GetResponseTimeTrends(&spanstore.ResponseTimeQueryParameters{
								BasicQueryParameters: spanstore.BasicQueryParameters{
									StartTimeMin: bqp.StartTimeMin,
									StartTimeMax: bqp.StartTimeMax,
								},
								ApplicationName: applicationName,
								Instance:        node.Name,
								TimeInterval:    time.Minute,
							})
						}
						return nil, nil
					},
				},
			},
		},
	)
}

func makeServerList(handler *APIHandler) *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(makeGLNode(handler)),
		Args: graphql.FieldConfigArgument{
			"duration": &graphql.ArgumentConfig{
				Type: gl.GLDurationType,
			},
			"applicationName": &graphql.ArgumentConfig{
				Type:         graphql.String,
				DefaultValue: "",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var durationParams gl.Duration
			err := mapstructure.Decode(p.Args["duration"], &durationParams)
			if err != nil {
				return nil, err
			}
			applicationName, ok := p.Args["applicationName"].(string)
			if !ok {
				return nil, errors.New("applicationName is not a valid Int")
			}
			extReader := handler.spanReader.(spanstore.ExtReader)
			bqp, err := durationParams.ToBasicQueryParameters()
			if err != nil {
				return nil, err
			}
			td, err := extReader.GetNodes(&spanstore.NodesQueryParameters{
				BasicQueryParameters: spanstore.BasicQueryParameters{
					StartTimeMin: bqp.StartTimeMin,
					StartTimeMax: bqp.StartTimeMax,
				},
				ApplicationName: applicationName,
			})
			retMe := make([]gl.Node, len(td))
			for i, t := range td {
				retMe[i].ApplicationName = applicationName
				retMe[i].Name = t
			}
			return retMe, err
		},
	}
}

func makeServiceThroughput(handler *APIHandler) *graphql.Field {
	return &graphql.Field{
		Type: gl.GLTrendListType,
		Args: graphql.FieldConfigArgument{
			"serviceName": &graphql.ArgumentConfig{
				Type:         graphql.ID,
				DefaultValue: "",
			},
			"duration": &graphql.ArgumentConfig{
				Type: gl.GLDurationType,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var durationParams gl.Duration
			err := mapstructure.Decode(p.Args["duration"], &durationParams)
			if err != nil {
				return nil, err
			}
			bqp, err := durationParams.ToBasicQueryParameters()
			if err != nil {
				return nil, err
			}
			extReader := handler.spanReader.(spanstore.ExtReader)
			if serviceName, ok := p.Args["serviceName"].(string); ok {
				ts, err := extReader.GetThroughputTrends(&spanstore.ThroughputQueryParameters{
					BasicQueryParameters: spanstore.BasicQueryParameters{
						StartTimeMin: bqp.StartTimeMin,
						StartTimeMax: bqp.StartTimeMax,
					},
					OperationName: serviceName,
					TimeInterval:  time.Minute,
				})
				if err != nil {
					return nil, err
				}
				return gl.Trends{TrendList: ts}, nil
			}
			return nil, errors.New("Invalid serviceName")
		},
	}
}

func makeServiceResponseTime(handler *APIHandler) *graphql.Field {
	return &graphql.Field{
		Type: gl.GLTrendListType,
		Args: graphql.FieldConfigArgument{
			"serviceName": &graphql.ArgumentConfig{
				Type: graphql.ID,
			},
			"duration": &graphql.ArgumentConfig{
				Type: gl.GLDurationType,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var durationParams gl.Duration
			err := mapstructure.Decode(p.Args["duration"], &durationParams)
			if err != nil {
				return nil, err
			}
			bqp, err := durationParams.ToBasicQueryParameters()
			if err != nil {
				return nil, err
			}
			extReader := handler.spanReader.(spanstore.ExtReader)
			if serviceName, ok := p.Args["serviceName"].(string); ok {
				ts, err := extReader.GetResponseTimeTrends(&spanstore.ResponseTimeQueryParameters{
					BasicQueryParameters: spanstore.BasicQueryParameters{
						StartTimeMin: bqp.StartTimeMin,
						StartTimeMax: bqp.StartTimeMax,
					},
					OperationName: serviceName,
					TimeInterval:  time.Minute,
				})
				if err != nil {
					return nil, err
				}
				retMe := gl.Trends{}
				retMe.TrendList = make([]int, len(ts))
				for i, t := range ts {
					retMe.TrendList[i] = int(t)
				}
				return retMe, nil
			}
			return nil, errors.New("Invalid serviceName")
		},
	}
}

func makeTraceList(handler *APIHandler) *graphql.Field {
	return &graphql.Field{
		Type: gl.GLTraceListType,
		Args: graphql.FieldConfigArgument{
			"condition": &graphql.ArgumentConfig{
				//DefaultValue:
				//Description:
				Type: gl.GLTraceQueryConditionType,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var condition gl.TraceQueryCondition
			err := mapstructure.Decode(p.Args["condition"], &condition)
			params, err := condition.ToTraceQueryParameters()
			if err != nil {
				return nil, err
			}
			traces, err := handler.spanReader.FindTraces(params)
			if err != nil {
				return nil, err
			}
			uiTraces := make([]*ui.Trace, len(traces))
			for i, v := range traces {
				uiTrace, uiErr := handler.convertModelToUI(v, true)
				if uiErr != nil {
					continue
				}
				uiTraces[i] = uiTrace
			}

			return gl.TraceList{
				Total:  len(traces),
				Traces: uiTraces,
			}, nil
		},
	}
}

func makeTrace(handler *APIHandler) *graphql.Field {
	return &graphql.Field{
		Type: gl.GLTraceType,
		Args: graphql.FieldConfigArgument{
			"traceId": &graphql.ArgumentConfig{
				//DefaultValue:
				//Description:
				Type: graphql.ID,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if traceId, ok := p.Args["traceId"].(string); ok {
				modeTraceId, err := model.TraceIDFromString(traceId)
				if err != nil {
					return nil, errors.New("Invalid traceId")
				}
				trace, err := handler.spanReader.GetTrace(modeTraceId)
				uiTrace, uiErr := handler.convertModelToUI(trace, true)
				if uiErr != nil {
					return nil, errors.New(uiErr.Msg)
				}
				return uiTrace, nil
			}
			return nil, errors.New("Invalid traceId")
		},
	}
}
