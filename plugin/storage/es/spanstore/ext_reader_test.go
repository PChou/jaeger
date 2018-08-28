package spanstore

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jaegertracing/jaeger/pkg/es"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gopkg.in/olivere/elastic.v5"
)

func getLocalReader() (*SpanReader, error) {
	innerClient, err := elastic.NewClient(elastic.SetURL("http://192.168.21.12:9200"))
	if err != nil {
		return nil, err
	}
	bulk, err := innerClient.BulkProcessor().Do(context.Background())
	if err != nil {
		return nil, err
	}
	client := es.WrapESClient(innerClient, bulk)
	return newSpanReader(client, zap.NewNop(), 0, ""), nil
}

func TestGetApplications(t *testing.T) {
	reader, err := getLocalReader()
	assert.Nil(t, err)

	ret, err := reader.GetApplications(&spanstore.BasicQueryParameters{
		StartTimeMin: time.Unix(1535339000, 0),
		StartTimeMax: time.Unix(1535342600, 0),
	})
	assert.Nil(t, err)
	fmt.Println(ret)
}

func TestGetLayerServices1(t *testing.T) {
	reader, err := getLocalReader()
	assert.Nil(t, err)

	ret, err := reader.GetLayerServices(&spanstore.LayerTypeQueryParameters{
		BasicQueryParameters: spanstore.BasicQueryParameters{
			StartTimeMin: time.Unix(1535339000, 0),
			StartTimeMax: time.Unix(1535342600, 0),
		},
		Layer: "cache",
		Type:  "exit",
		By:    tagPeerKeyField,
	})
	assert.Nil(t, err)
	fmt.Println(ret)
}

func TestGetLayerServices2(t *testing.T) {
	reader, err := getLocalReader()
	assert.Nil(t, err)

	ret, err := reader.GetLayerServices(&spanstore.LayerTypeQueryParameters{
		BasicQueryParameters: spanstore.BasicQueryParameters{
			StartTimeMin: time.Unix(1535339000, 0),
			StartTimeMax: time.Unix(1535342600, 0),
		},
		Layer: "http",
		Type:  "entry",
		By:    operationNameField,
	})
	assert.Nil(t, err)
	fmt.Println(ret)
}

func TestGetThermoDynamic(t *testing.T) {
	reader, err := getLocalReader()
	assert.Nil(t, err)

	ret, err := reader.GetThermoDynamic(&spanstore.ThermoDynamicQueryParameters{
		BasicQueryParameters: spanstore.BasicQueryParameters{
			StartTimeMin: time.Unix(1535339000, 0),
			StartTimeMax: time.Unix(1535342600, 0),
		},
		TimeInterval:            time.Minute,
		DurationInterval:        time.Millisecond * 100,
		DurationExtendBoundsMin: 0,
		DurationExtendBoundsMax: time.Millisecond * 3000,
	})
	assert.Nil(t, err)
	fmt.Println(ret)
}

func TestGetServiceTopResponseTime(t *testing.T) {
	reader, err := getLocalReader()
	assert.Nil(t, err)
	ret, err := reader.GetServiceTopResponseTime(&spanstore.ServiceTopResponseTimeQueryParameters{
		BasicQueryParameters: spanstore.BasicQueryParameters{
			StartTimeMin: time.Unix(1535339000, 0),
			StartTimeMax: time.Unix(1535342600, 0),
		},
		Top: 3,
	})
	assert.Nil(t, err)
	fmt.Println(ret)
}

func TestGetServiceThroughput(t *testing.T) {
	reader, err := getLocalReader()
	assert.Nil(t, err)
	ret, err := reader.GetServiceThroughput(&spanstore.ThroughputQueryParameters{
		BasicQueryParameters: spanstore.BasicQueryParameters{
			StartTimeMin: time.Unix(1535339000, 0),
			StartTimeMax: time.Unix(1535342600, 0),
		},
		OperationName: "/sample",
		TimeInterval:  time.Minute,
	})
	assert.Nil(t, err)
	fmt.Println(ret)
}

func TestGetServiceResponseTime(t *testing.T) {
	reader, err := getLocalReader()
	assert.Nil(t, err)
	ret, err := reader.GetServiceResponseTime(&spanstore.ResponseTimeQueryParameters{
		BasicQueryParameters: spanstore.BasicQueryParameters{
			StartTimeMin: time.Unix(1535339000, 0),
			StartTimeMax: time.Unix(1535342600, 0),
		},
		OperationName: "/sample",
		TimeInterval:  time.Minute,
	})
	assert.Nil(t, err)
	fmt.Println(ret)
}

func TestGetNodeThroughput(t *testing.T) {
	reader, err := getLocalReader()
	assert.Nil(t, err)
	ret, err := reader.GetNodeThroughput(&spanstore.ThroughputQueryParameters{
		BasicQueryParameters: spanstore.BasicQueryParameters{
			StartTimeMin: time.Unix(1535356800, 0),
			StartTimeMax: time.Unix(1535360400, 0),
		},
		//ApplicationName: "jboss2",
		TimeInterval: time.Minute,
	})
	assert.Nil(t, err)
	for _, r := range ret {
		fmt.Println(r)
	}
}

func TestGetApplicationTopThroughput(t *testing.T) {
	reader, err := getLocalReader()
	assert.Nil(t, err)
	ret, err := reader.GetApplicationTopThroughput(&spanstore.TopThroughputQueryParameters{
		BasicQueryParameters: spanstore.BasicQueryParameters{
			StartTimeMin: time.Unix(1535356800, 0),
			StartTimeMax: time.Unix(1535360400, 0),
		},
	})
	assert.Nil(t, err)
	for _, r := range ret {
		fmt.Println(r)
	}
}

func TestGetNodeTopThroughput(t *testing.T) {
	reader, err := getLocalReader()
	assert.Nil(t, err)
	ret, err := reader.GetNodeTopThroughput(&spanstore.TopThroughputQueryParameters{
		BasicQueryParameters: spanstore.BasicQueryParameters{
			StartTimeMin: time.Unix(1535356800, 0),
			StartTimeMax: time.Unix(1535360400, 0),
		},
	})
	assert.Nil(t, err)
	for _, r := range ret {
		fmt.Println(r)
	}
}
