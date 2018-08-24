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
	GetApplications(query *ApplicationQueryParameter) ([]string, error)
	GetThermoDynamic(query *ThermoDynamicQueryParameters) (*model.ThermoDynamic, error)
	GetTrends(query *TrendsQueryParameters) ([]int, error)
}

type ThermoDynamicQueryParameters struct {
	ServiceName             string
	OperationName           string
	StartTimeMin            time.Time
	StartTimeMax            time.Time
	TimeInterval            time.Duration
	DurationInterval        time.Duration
	DurationExtendBoundsMin time.Duration
	DurationExtendBoundsMax time.Duration
}

type TrendsQueryParameters struct {
	ServiceName   string
	OperationName string
	StartTimeMin  time.Time
	StartTimeMax  time.Time
	TimeInterval  time.Duration
}

type ApplicationQueryParameter struct {
	StartTimeMin time.Time
	StartTimeMax time.Time
}
