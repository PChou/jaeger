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
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"

	"github.com/Shopify/sarama"
	ss "github.com/jaegertracing/jaeger/cmd/collector/app/sampling/strategystore"
	"github.com/jaegertracing/jaeger/thrift-gen/sampling"
	"github.com/pkg/errors"
)

type strategyStore struct {
	logger  *zap.Logger
	metrics *kafkaStrategyStoreMetrics

	brokers []string
	topic   string
	exit    chan os.Signal
	lock    *sync.RWMutex //the lock used to protect the serviceStrategies map

	defaultStrategy   *sampling.SamplingStrategyResponse
	serviceStrategies map[string]*sampling.SamplingStrategyResponse
}

// NewStrategyStore creates a strategy store that holds static sampling strategies.
func NewStrategyStore(options Options, logger *zap.Logger) (ss.StrategyStore, error) {
	h := &strategyStore{
		brokers:           options.Brokers,
		topic:             options.Topic,
		logger:            logger,
		metrics:           newMetrics(options),
		exit:              make(chan os.Signal, 0),
		lock:              new(sync.RWMutex),
		defaultStrategy:   &defaultStrategy,
		serviceStrategies: make(map[string]*sampling.SamplingStrategyResponse),
	}

	//start a seperate goroutine subscribe from kafka
	config := sarama.NewConfig()
	config.Version = sarama.V0_10_0_0
	consumer, err := sarama.NewConsumer(h.brokers, config)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to create consumer to brokers %v", h.brokers))
	}
	ps, err := consumer.Partitions(h.topic)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to fetch partition info for the topic %s", h.topic))
	}
	signal.Notify(h.exit, os.Interrupt, syscall.SIGTERM)
	h.metrics.ConnectionState.Update(1)
	for _, partition := range ps {
		pConsumer, err := consumer.ConsumePartition(h.topic, partition, sarama.OffsetOldest)
		if err != nil {
			continue
		}
		go func() {
			for {
				select {
				case <-h.exit:
					pConsumer.Close()
					return
				case message := <-pConsumer.Messages():
					h.metrics.ConsumeCounter.Inc(1)
					h.apply(message)
				}
			}
		}()
	}

	return h, nil
}

func (h *strategyStore) apply(message *sarama.ConsumerMessage) {
	if message == nil {
		return
	}
	var s strategy
	err := json.Unmarshal(message.Value, &s)
	if err != nil {
		h.logger.Error("Failed to unmarshal kafka message")
		return
	}
	h.lock.Lock()
	defer h.lock.Unlock()
	h.metrics.ApplyCounter.Inc(1)
	// nil Key means default strategy
	if message.Key == nil {
		h.set(h.defaultStrategy, s)
	} else {
		if _, ok := h.serviceStrategies[string(message.Key)]; !ok {
			h.serviceStrategies[string(message.Key)] = &sampling.SamplingStrategyResponse{}
		}
		h.set(h.serviceStrategies[string(message.Key)], s)
	}
}

func (h *strategyStore) set(t *sampling.SamplingStrategyResponse, s strategy) {
	if s.Type == samplerTypeProbabilistic {
		t.StrategyType = sampling.SamplingStrategyType_PROBABILISTIC
		if t.ProbabilisticSampling == nil {
			t.ProbabilisticSampling = &sampling.ProbabilisticSamplingStrategy{}
		}
		t.ProbabilisticSampling.SamplingRate = s.Param
	} else {
		t.StrategyType = sampling.SamplingStrategyType_RATE_LIMITING
		if t.RateLimitingSampling == nil {
			t.RateLimitingSampling = &sampling.RateLimitingSamplingStrategy{}
		}
		t.RateLimitingSampling.MaxTracesPerSecond = int16(s.Param)
	}
}

// GetSamplingStrategy implements StrategyStore#GetSamplingStrategy.
func (h *strategyStore) GetSamplingStrategy(serviceName string) (*sampling.SamplingStrategyResponse, error) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	//h.logger.Info("getting sampling strategy", zap.Field{Key: "service_name", Type: zapcore.StringType, String: serviceName})
	h.metrics.getServiceCounter(serviceName).Inc(1)
	if strategy, ok := h.serviceStrategies[serviceName]; ok {
		return strategy, nil
	}
	return h.defaultStrategy, nil
}
