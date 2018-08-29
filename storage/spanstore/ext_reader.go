// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spanstore

import (
	"time"

	"github.com/jaegertracing/jaeger/model"
)

// ExtReader used to support more complex query, like thermodynamic graph
type ExtReader interface {
	GetApplications(query *BasicQueryParameters) ([]string, error)
	GetLayerServices(query *LayerTypeQueryParameters) ([]string, error)
	GetNodes(query *NodesQueryParameters) ([]string, error)
	GetServiceTopResponseTime(query *ServiceTopResponseTimeQueryParameters) ([]*model.ServiceAvgResponseTime, error)
	GetThermoDynamic(query *ThermoDynamicQueryParameters) (*model.ThermoDynamic, error)
	GetApplicationTopThroughput(query *TopThroughputQueryParameters) ([]*model.ApplicationThroughput, error)
	GetNodeTopThroughput(query *TopThroughputQueryParameters) ([]*model.NodeAvgThroughput, error)
	GetResponseTimeTrends(query *ResponseTimeQueryParameters) ([]float64, error)
	GetThroughputTrends(query *ThroughputQueryParameters) ([]int, error)
}

// Basic query parameters contains StartTimeMin and StartTimeMax
// Almost every query should contains the StartTime duration
type BasicQueryParameters struct {
	StartTimeMin time.Time
	StartTimeMax time.Time
}

type LayerTypeQueryParameters struct {
	BasicQueryParameters
	ApplicationName string
	Layer           string
	Type            string
	By              string
}

type ServiceTopResponseTimeQueryParameters struct {
	BasicQueryParameters
	ApplicationName string
	Top             int
}

type ThermoDynamicQueryParameters struct {
	BasicQueryParameters
	TimeInterval            time.Duration
	DurationInterval        time.Duration
	DurationExtendBoundsMin time.Duration
	DurationExtendBoundsMax time.Duration
}

type TopThroughputQueryParameters struct {
	BasicQueryParameters
	ApplicationName string
	Top             int
}

type ThroughputQueryParameters struct {
	BasicQueryParameters
	ApplicationName string
	OperationName   string
	Instance        string
	TimeInterval    time.Duration
}

type ResponseTimeQueryParameters struct {
	BasicQueryParameters
	ApplicationName string
	OperationName   string
	Instance        string
	TimeInterval    time.Duration
}

type NodesQueryParameters struct {
	BasicQueryParameters
	ApplicationName string
	OperationName   string
}
