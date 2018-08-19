// Copyright (c) 2018 The Jaeger Authors.
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

package kafka

import (
	"strings"

	"github.com/uber/jaeger-lib/metrics"
)

// Metrics set
type kafkaStrategyStoreMetrics struct {
	factory         metrics.Factory
	ConnectionState metrics.Gauge
	ConsumeCounter  metrics.Counter
	ApplyCounter    metrics.Counter
	ServiceCounter  map[string]metrics.Counter
}

func newMetrics(option Options) *kafkaStrategyStoreMetrics {
	return &kafkaStrategyStoreMetrics{
		factory:         option.metricsFactory,
		ConnectionState: option.metricsFactory.Gauge("connection", map[string]string{"brokers": strings.Join(option.Brokers, ",")}),
		ConsumeCounter:  option.metricsFactory.Counter("consume", map[string]string{"topic": option.Topic}),
		ApplyCounter:    option.metricsFactory.Counter("apply", map[string]string{"topic": option.Topic}),
		ServiceCounter:  make(map[string]metrics.Counter),
	}
}

func (m *kafkaStrategyStoreMetrics) getServiceCounter(name string) metrics.Counter {
	if _, ok := m.ServiceCounter[name]; !ok {
		m.ServiceCounter[name] = m.factory.Counter("request", map[string]string{"service_name": name})
	}
	return m.ServiceCounter[name]
}
