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
	"flag"
	"strings"

	"github.com/spf13/viper"
	"github.com/uber/jaeger-lib/metrics"
)

const (
	// Brokers list of a kafka cluster
	samplingStrategiesKafka = "sampling.strategies-kafka"
	// topic used, to fetch sampling strategies from
	samplingStrategiesKafkaTopic = "sampling.strategies-kafka-topic"
)

// Options holds configuration for the kafka sampling strategy store.
type Options struct {
	// Brokers is the brokers list of a kafka cluster
	Brokers []string
	// Topic used, to fetch sampling strategies from
	Topic string

	metricsFactory metrics.Factory
}

// AddFlags adds flags for Options
func AddFlags(flagSet *flag.FlagSet) {
	flagSet.String(samplingStrategiesKafka, "", "The broker list of kafka, host1:port1,host2:port2")
	flagSet.String(samplingStrategiesKafkaTopic, "jaeger-sampling-strategies", "The topic subscribed from, to change the sampling strategies of default or each service")
}

// InitFromViper initializes Options with properties from viper
func (opts *Options) InitFromViper(v *viper.Viper) *Options {
	// TODO trim
	opts.Brokers = strings.Split(v.GetString(samplingStrategiesKafka), ",")
	opts.Topic = v.GetString(samplingStrategiesKafkaTopic)
	return opts
}
