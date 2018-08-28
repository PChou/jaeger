package spanstore

import (
	"fmt"
	"time"

	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	"github.com/pkg/errors"
	"gopkg.in/olivere/elastic.v5"
)

const (
	startTimeMillisField = "startTimeMillis"
	//Service Count: span.layer =
	//DB Count: span.layer = db and group by peer
	//Cache Count: span.layer = cache and group by peer
	tagLayerKeyField    = "flattenTags.span.layer"
	tagTypeKeyField     = "flattenTags.span.type"
	tagPeerKeyField     = "flattenTags.peer"
	tagInstanceKeyField = "flattenTags.process.sid"

	serviceLayer   = "HTTP"
	serviceDefType = "Entry"
	dbLayer        = "DB"
	dbDefType      = "Exit"
	cacheLayer     = "CACHE"
	cacheDefType   = "Exit"
)

func (s *SpanReader) buildStartTimeMillisQuery(startTimeMin time.Time, startTimeMax time.Time) elastic.Query {
	minStartTimeMilli := model.TimeAsEpochMilliseconds(startTimeMin)
	maxStartTimeMilli := model.TimeAsEpochMilliseconds(startTimeMax)
	return elastic.NewRangeQuery(startTimeMillisField).Gte(minStartTimeMilli).Lte(maxStartTimeMilli)
}

func (s *SpanReader) GetApplications(query *spanstore.BasicQueryParameters) ([]string, error) {
	serviceIndices := s.indicesForTimeRange(s.serviceIndexPrefix, query.StartTimeMin, query.StartTimeMax)
	return s.serviceOperationStorage.getServices(serviceIndices)
}

func (s *SpanReader) GetLayerServices(query *spanstore.LayerTypeQueryParameters) ([]string, error) {
	timeRange := s.buildStartTimeMillisQuery(query.StartTimeMin, query.StartTimeMax)
	whereQuery := elastic.NewBoolQuery().Must(timeRange)
	if query.Layer != "" {
		whereQuery.Must(elastic.NewMatchQuery(tagLayerKeyField, query.Layer))
	}
	if query.Type != "" {
		whereQuery.Must(elastic.NewMatchQuery(tagTypeKeyField, query.Type))
	}
	if query.ApplicationName != "" {
		whereQuery.Must(elastic.NewMatchQuery(serviceNameField, query.ApplicationName))
	}

	var agg *elastic.TermsAggregation
	if query.By == operationNameField {
		agg = elastic.NewTermsAggregation().Field(operationNameField)
	} else if query.By == tagPeerKeyField {
		agg = elastic.NewTermsAggregation().Field(tagPeerKeyField)
	}
	jaegerIndices := s.indicesForTimeRange(s.spanIndexPrefix, query.StartTimeMin, query.StartTimeMax)
	searchService := s.client.Search(jaegerIndices...).
		Type(spanType).
		Size(0). // set to 0 because we don't want actual documents.
		Aggregation("agg", agg).
		IgnoreUnavailable(true).
		Query(whereQuery)

	searchResult, err := searchService.Do(s.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Search service failed")
	}

	bucket, found := searchResult.Aggregations.Terms("agg")
	if !found {
		return nil, errors.New("Counld not found bucket by aggs")
	}

	retMe := make([]string, len(bucket.Buckets))
	for i, b := range bucket.Buckets {
		retMe[i] = fmt.Sprintf("%v", b.Key)
	}
	// b, _ := json.Marshal(searchResult)
	// fmt.Println(string(b))
	return retMe, nil
}

func (s *SpanReader) GetServiceTopResponseTime(query *spanstore.ServiceTopResponseTimeQueryParameters) ([]*model.ServiceAvgResponseTime, error) {
	timeRange := s.buildStartTimeMillisQuery(query.StartTimeMin, query.StartTimeMax)
	whereQuery := elastic.NewBoolQuery().Must(timeRange)
	whereQuery.Must(elastic.NewMatchQuery(tagLayerKeyField, serviceLayer))
	whereQuery.Must(elastic.NewMatchQuery(tagTypeKeyField, serviceDefType))
	if query.ApplicationName != "" {
		whereQuery.Must(elastic.NewMatchQuery(serviceNameField, query.ApplicationName))
	}

	agg := elastic.NewTermsAggregation().Field(operationNameField).
		SubAggregation("avg", elastic.NewAvgAggregation().Field(durationField)).
		//SubAggregation("avg", elastic.NewAvgAggregation().Field(durationField)).
		OrderByAggregation("avg", false)

	jaegerIndices := s.indicesForTimeRange(s.spanIndexPrefix, query.StartTimeMin, query.StartTimeMax)
	searchService := s.client.Search(jaegerIndices...).
		Type(spanType).
		Size(0). // set to 0 because we don't want actual documents.
		Aggregation("agg", agg).
		IgnoreUnavailable(true)
	searchResult, err := searchService.Do(s.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Search service failed")
	}

	opBucket, found := searchResult.Aggregations.Terms("agg")
	if !found {
		return nil, errors.New("Counld not found bucket by aggs")
	}

	retMe := make([]*model.ServiceAvgResponseTime, 0, len(opBucket.Buckets))
	for i, b := range opBucket.Buckets {
		if i >= query.Top {
			break
		}
		if b.DocCount == 0 {
			retMe = append(retMe, &model.ServiceAvgResponseTime{})
		} else {
			avg, found := b.Avg("avg")
			if !found {
				retMe = append(retMe, &model.ServiceAvgResponseTime{ServiceName: fmt.Sprintf("%v", b.Key), Value: 0})
			} else {
				retMe = append(retMe, &model.ServiceAvgResponseTime{ServiceName: fmt.Sprintf("%v", b.Key), Value: *avg.Value})
			}
		}
	}

	// b, _ := json.Marshal(searchResult)
	// fmt.Println(string(b))
	return retMe, nil

}

func (s *SpanReader) GetThermoDynamic(query *spanstore.ThermoDynamicQueryParameters) (*model.ThermoDynamic, error) {
	// {
	// 	"aggregations": {
	// 		"date_histogram": {
	// 			"aggregations": {
	// 				"histogram": {
	// 					"histogram": {
	// 						"extended_bounds": {
	// 							"max": 3000000,
	// 							"min": 0
	// 						},
	// 						"field": "duration",
	// 						"interval": 200000
	// 					}
	// 				}
	// 			},
	// 			"histogram": {
	// 				"extended_bounds": {
	// 					"max": 1534944600000,
	// 					"min": 1534943700000
	// 				},
	// 				"field": "startTimeMillis",
	// 				"interval": 60000
	// 			}
	// 		}
	// 	},
	// 	"query": {
	// 		"bool": {
	// 			"must": {
	// 				"range": {
	// 					"startTimeMillis": {
	// 						"from": 1534943700000,
	// 						"include_lower": true,
	// 						"include_upper": true,
	// 						"to": 1534944600000
	// 					}
	// 				}
	// 			}
	//
	// 	},
	// 	"size": 0
	// }

	//TODO: check the validation of the query
	timeRange := s.buildStartTimeMillisQuery(query.StartTimeMin, query.StartTimeMax)
	whereQuery := elastic.NewBoolQuery().Must(timeRange)
	minStartTimeMilli := model.TimeAsEpochMilliseconds(query.StartTimeMin)
	maxStartTimeMilli := model.TimeAsEpochMilliseconds(query.StartTimeMax)
	durationAgg := elastic.NewHistogramAggregation().
		Field(durationField).
		Interval(float64(query.DurationInterval/time.Microsecond)). //duration is microseconds
		ExtendedBounds(float64(query.DurationExtendBoundsMin/time.Microsecond), float64(query.DurationExtendBoundsMax/time.Microsecond))

	timeAgg := elastic.NewHistogramAggregation().
		Field(startTimeMillisField).
		Interval(float64(query.TimeInterval/time.Millisecond)).
		ExtendedBounds(float64(minStartTimeMilli), float64(maxStartTimeMilli)).
		SubAggregation("histogram", durationAgg)

	jaegerIndices := s.indicesForTimeRange(s.spanIndexPrefix, query.StartTimeMin, query.StartTimeMax)
	searchService := s.client.Search(jaegerIndices...).
		Type(spanType).
		Size(0). // set to 0 because we don't want actual documents.
		Aggregation("date_histogram", timeAgg).
		IgnoreUnavailable(true).
		Query(whereQuery)

	searchResult, err := searchService.Do(s.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Search service failed")
	}

	timeBulket, found := searchResult.Aggregations.Terms("date_histogram")
	if !found {
		return nil, errors.New("Counld not found bucket by date_histogram")
	}

	retMe := &model.ThermoDynamic{}
	retMe.ResponseTimeStep = int(query.DurationInterval / time.Millisecond)
	retMe.Nodes = make([][3]int, 0)
	for i, dateBulk := range timeBulket.Buckets {
		durationBulket, found := dateBulk.Terms("histogram")
		if !found {
			continue
		}
		column := make([][3]int, len(durationBulket.Buckets))
		for j, durationBulk := range durationBulket.Buckets {
			column[j][0] = i
			column[j][1] = j
			column[j][2] = int(durationBulk.DocCount)
		}
		retMe.Nodes = append(retMe.Nodes, column...)
	}
	//b, _ := json.Marshal(searchResult)
	//fmt.Println(string(b))
	return retMe, nil
}

func (s *SpanReader) GetResponseTimeTrends(query *spanstore.ResponseTimeQueryParameters) ([]float64, error) {
	timeRange := s.buildStartTimeMillisQuery(query.StartTimeMin, query.StartTimeMax)
	whereQuery := elastic.NewBoolQuery().Must(timeRange)
	if query.ApplicationName != "" {
		whereQuery.Must(elastic.NewMatchQuery(serviceNameField, query.ApplicationName))
	}
	if query.OperationName != "" {
		whereQuery.Must(elastic.NewMatchQuery(operationNameField, query.OperationName))
	}
	if query.Instance != "" {
		whereQuery.Must(elastic.NewMatchQuery(tagInstanceKeyField, query.Instance))
	}
	avgAgg := elastic.NewAvgAggregation().Field(durationField)
	minStartTimeMilli := model.TimeAsEpochMilliseconds(query.StartTimeMin)
	maxStartTimeMilli := model.TimeAsEpochMilliseconds(query.StartTimeMax)
	timeAgg := elastic.NewHistogramAggregation().
		Field(startTimeMillisField).
		Interval(float64(query.TimeInterval/time.Millisecond)).
		ExtendedBounds(float64(minStartTimeMilli), float64(maxStartTimeMilli)).
		SubAggregation("duration_avg", avgAgg)

	jaegerIndices := s.indicesForTimeRange(s.spanIndexPrefix, query.StartTimeMin, query.StartTimeMax)
	searchService := s.client.Search(jaegerIndices...).
		Type(spanType).
		Size(0). // set to 0 because we don't want actual documents.
		Aggregation("date_histogram", timeAgg).
		IgnoreUnavailable(true).
		Query(whereQuery)

	searchResult, err := searchService.Do(s.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Search service failed")
	}
	timeBulket, found := searchResult.Aggregations.Terms("date_histogram")
	if !found {
		return nil, errors.New("Counld not found bucket by date_histogram")
	}

	retMe := make([]float64, len(timeBulket.Buckets))
	for i, b := range timeBulket.Buckets {
		if b.DocCount == 0 {
			retMe[i] = 0
		} else {
			avg, found := b.Avg("duration_avg")
			if !found {
				retMe[i] = 0
			} else {
				retMe[i] = *avg.Value
			}
		}
	}
	// b, _ := json.Marshal(searchResult)
	// fmt.Println(string(b))
	return retMe, nil
}

func (s *SpanReader) GetApplicationTopThroughput(query *spanstore.TopThroughputQueryParameters) ([]*model.ApplicationThroughput, error) {
	timeRange := s.buildStartTimeMillisQuery(query.StartTimeMin, query.StartTimeMax)
	whereQuery := elastic.NewBoolQuery().Must(timeRange)

	appAgg := elastic.NewTermsAggregation().Field(serviceNameField)
	jaegerIndices := s.indicesForTimeRange(s.spanIndexPrefix, query.StartTimeMin, query.StartTimeMax)
	searchService := s.client.Search(jaegerIndices...).
		Type(spanType).
		Size(0). // set to 0 because we don't want actual documents.
		Aggregation("app", appAgg).
		IgnoreUnavailable(true).
		Query(whereQuery)
	searchResult, err := searchService.Do(s.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Search service failed")
	}
	appBulket, found := searchResult.Aggregations.Terms("app")
	if !found {
		return nil, errors.New("Counld not found bucket by app")
	}
	retMe := make([]*model.ApplicationThroughput, len(appBulket.Buckets))
	minutesCount := query.StartTimeMax.Sub(query.StartTimeMin) / time.Minute
	for i, b := range appBulket.Buckets {
		if b.DocCount == 0 {
			retMe[i] = &model.ApplicationThroughput{}
		} else {
			retMe[i] = &model.ApplicationThroughput{
				ApplicationName: fmt.Sprintf("%v", b.Key),
				Value:           float64(b.DocCount) / float64(minutesCount),
			}
		}
	}
	// b, _ := json.Marshal(searchResult)
	// fmt.Println(string(b))
	return retMe, nil
}

func (s *SpanReader) GetNodeTopThroughput(query *spanstore.TopThroughputQueryParameters) ([]*model.NodeAvgThroughput, error) {
	timeRange := s.buildStartTimeMillisQuery(query.StartTimeMin, query.StartTimeMax)
	whereQuery := elastic.NewBoolQuery().Must(timeRange)
	if query.ApplicationName != "" {
		whereQuery.Must(elastic.NewMatchQuery(serviceNameField, query.ApplicationName))
	}

	appAgg := elastic.NewTermsAggregation().Field(tagInstanceKeyField).OrderByCount(false)
	jaegerIndices := s.indicesForTimeRange(s.spanIndexPrefix, query.StartTimeMin, query.StartTimeMax)
	searchService := s.client.Search(jaegerIndices...).
		Type(spanType).
		Size(0). // set to 0 because we don't want actual documents.
		Aggregation("instance", appAgg).
		IgnoreUnavailable(true).
		Query(whereQuery)
	searchResult, err := searchService.Do(s.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Search service failed")
	}
	appBulket, found := searchResult.Aggregations.Terms("instance")
	if !found {
		return nil, errors.New("Counld not found bucket by instance")
	}
	retMe := make([]*model.NodeAvgThroughput, len(appBulket.Buckets))
	minutesCount := query.StartTimeMax.Sub(query.StartTimeMin) / time.Minute
	for i, b := range appBulket.Buckets {
		if b.DocCount == 0 {
			retMe[i] = &model.NodeAvgThroughput{}
		} else {
			retMe[i] = &model.NodeAvgThroughput{
				Node:  fmt.Sprintf("%v", b.Key),
				Value: float64(b.DocCount) / float64(minutesCount),
			}
		}
	}
	// b, _ := json.Marshal(searchResult)
	// fmt.Println(string(b))
	return retMe, nil
}

func (s *SpanReader) GetThroughputTrends(query *spanstore.ThroughputQueryParameters) ([]int, error) {
	timeRange := s.buildStartTimeMillisQuery(query.StartTimeMin, query.StartTimeMax)
	whereQuery := elastic.NewBoolQuery().Must(timeRange)
	if query.ApplicationName != "" {
		whereQuery.Must(elastic.NewMatchQuery(serviceNameField, query.ApplicationName))
	}
	if query.OperationName != "" {
		whereQuery.Must(elastic.NewMatchQuery(operationNameField, query.OperationName))
	}
	if query.Instance != "" {
		whereQuery.Must(elastic.NewMatchQuery(tagInstanceKeyField, query.Instance))
	}

	minStartTimeMilli := model.TimeAsEpochMilliseconds(query.StartTimeMin)
	maxStartTimeMilli := model.TimeAsEpochMilliseconds(query.StartTimeMax)

	timeAgg := elastic.NewHistogramAggregation().
		Field(startTimeMillisField).
		Interval(float64(query.TimeInterval/time.Millisecond)).
		ExtendedBounds(float64(minStartTimeMilli), float64(maxStartTimeMilli))

	jaegerIndices := s.indicesForTimeRange(s.spanIndexPrefix, query.StartTimeMin, query.StartTimeMax)
	searchService := s.client.Search(jaegerIndices...).
		Type(spanType).
		Size(0). // set to 0 because we don't want actual documents.
		Aggregation("date_histogram", timeAgg).
		IgnoreUnavailable(true).
		Query(whereQuery)

	searchResult, err := searchService.Do(s.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Search service failed")
	}
	sidBulket, found := searchResult.Aggregations.Terms("date_histogram")
	if !found {
		return nil, errors.New("Counld not found bucket by sid")
	}

	retMe := make([]int, len(sidBulket.Buckets))
	for i, b := range sidBulket.Buckets {
		retMe[i] = int(b.DocCount)
	}
	// b, _ := json.Marshal(searchResult)
	// fmt.Println(string(b))
	return retMe, nil
}

func (s *SpanReader) GetNodes(query *spanstore.NodesQueryParameters) ([]string, error) {
	timeRange := s.buildStartTimeMillisQuery(query.StartTimeMin, query.StartTimeMax)
	whereQuery := elastic.NewBoolQuery().Must(timeRange)
	if query.ApplicationName != "" {
		whereQuery.Must(elastic.NewMatchQuery(serviceNameField, query.ApplicationName))
	}
	if query.OperationName != "" {
		whereQuery.Must(elastic.NewMatchQuery(operationNameField, query.OperationName))
	}

	insAgg := elastic.NewTermsAggregation().Field(tagInstanceKeyField).OrderByCount(false)
	jaegerIndices := s.indicesForTimeRange(s.spanIndexPrefix, query.StartTimeMin, query.StartTimeMax)
	searchService := s.client.Search(jaegerIndices...).
		Type(spanType).
		Size(0). // set to 0 because we don't want actual documents.
		Aggregation("instance", insAgg).
		IgnoreUnavailable(true).
		Query(whereQuery)
	searchResult, err := searchService.Do(s.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Search service failed")
	}
	appBulket, found := searchResult.Aggregations.Terms("instance")
	if !found {
		return nil, errors.New("Counld not found bucket by instance")
	}
	retMe := make([]string, len(appBulket.Buckets))
	for i, b := range appBulket.Buckets {
		retMe[i] = fmt.Sprintf("%v", b.Key)
	}
	// b, _ := json.Marshal(searchResult)
	// fmt.Println(string(b))
	return retMe, nil
}
