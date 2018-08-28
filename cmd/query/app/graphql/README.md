# graphql 文档生成

```sh
graphdoc -e http://localhost:16686/api/graphql -o schema
```

https://github.com/2fd/graphdoc

# 适配skywalking的GraphQL接口说明

实际处理请求的只有一个handler，但是在接口暴露上是分不同路由的。换句话说，下面这些路由都是可以用于所有的查询语句的：

```
/api/graphql
/api/trace
/api/trace/options
/api/spans
/api/dashboard
/api/service
```

# 查询Application

```json
{
	"query": "query q($duration: Duration) { applicationList(duration: $duration) { apps { name } count }}",
	"variables": { "duration":{ "start": "2018-08-22 13:15", "end": "2018-08-22 13:30", "step": "MUNITE"  }  }
}
```

在这个查询中`start`和`end`用于定位命中的索引`jaeger-service-xxxx`

返回样例：

```json
{"data":{"applicationList":{"apps":[{"name":"gateway"},{"name":"client"},{"name":"query-service"},{"name":"transfer-service"},{"name":"jaeger-query"}],"count":5}}}
```

# 查询热力图

```json
{
    "query": "query q($duration: Duration) { thermodynamic(duration: $duration) { responseTimeStep nodes  } }",
            "variables": { "duration":{ "start": "2018-08-22 13:15", "end": "2018-08-22 13:30", "step": "MUNITE"  }  }
}
```

在这个查询中`start`和`end`用于定位`jaeger-span-xxxx`的索引，以及过滤在这个时间范围内的span，`step`只支持`MUNITE`。

二级分桶（即根据span的duration分桶），程序固定用的是500ms，上下界为0~3000ms。上下界的概念详见es：[Histogram Aggregation](https://www.elastic.co/guide/en/elasticsearch/reference/current/search-aggregations-bucket-histogram-aggregation.html)


```go
return &spanstore.ThermoDynamicQueryParameters{
    StartTimeMin:            start,
    StartTimeMax:            end,
    TimeInterval:            time.Minute,
    DurationInterval:        time.Millisecond * 500,
    DurationExtendBoundsMin: 0,
    DurationExtendBoundsMax: time.Millisecond * 3000,
}, nil
```

样例返回结果

```json
{"data":{"thermodynamic":{"nodes":[[0,0,0],[0,1,0],[0,2,0],[0,3,0],[0,4,0],[0,5,0],[0,6,0],[1,0,0],[1,1,0],[1,2,0],[1,3,0],[1,4,0],[1,5,0],[1,6,0],[2,0,0],[2,1,0],[2,2,0],[2,3,0],[2,4,0],[2,5,0],[2,6,0],[3,0,0],[3,1,0],[3,2,0],[3,3,0],[3,4,0],[3,5,0],[3,6,0],[4,0,31],[4,1,15],[4,2,3],[4,3,0],[4,4,0],[4,5,0],[4,6,0],[5,0,40],[5,1,12],[5,2,3],[5,3,6],[5,4,2],[5,5,0],[5,6,0],[6,0,36],[6,1,13],[6,2,1],[6,3,3],[6,4,2],[6,5,0],[6,6,0],[7,0,36],[7,1,17],[7,2,5],[7,3,2],[7,4,4],[7,5,0],[7,6,0],[8,0,36],[8,1,16],[8,2,3],[8,3,2],[8,4,4],[8,5,0],[8,6,0],[9,0,28],[9,1,14],[9,2,10],[9,3,4],[9,4,6],[9,5,0],[9,6,0],[10,0,43],[10,1,9],[10,2,4],[10,3,5],[10,4,2],[10,5,4],[10,6,0],[11,0,55],[11,1,4],[11,2,5],[11,3,7],[11,4,0],[11,5,2],[11,6,0],[12,0,33],[12,1,15],[12,2,6],[12,3,6],[12,4,2],[12,5,2],[12,6,0],[13,0,14],[13,1,16],[13,2,2],[13,3,0],[13,4,0],[13,5,0],[13,6,0],[14,0,0],[14,1,0],[14,2,0],[14,3,0],[14,4,0],[14,5,0],[14,6,0],[15,0,0],[15,1,0],[15,2,0],[15,3,0],[15,4,0],[15,5,0],[15,6,0]],"responseTimeStep":500}}}
```


# 查询trace列表

```json
{
	"query": "query test($condition: TraceQueryCondition) { traceList(condition: $condition) { total traces { traceID spans { operationName startTime tags { key type value } } } }}",
	"variables": {"condition":{"minTraceDuration":0,"maxTraceDuration":3000,"traceState":"ALL","queryOrder":"BY_START_TIME","queryDuration":{"start":"2018-08-22 13:15","end":"2018-08-22 13:30","step":"MINUTE"},"paging":{"pageNum":1,"pageSize":20,"needTotal":true} } }
}
```

在这个查询中, GraphQL查询用的selection的层级写法可以参考jaeger原有的接口：http://192.168.32.29:16686/api/traces/77f1552cecc4c5a5

> 注意：该接口实现目前并没有做分页支持，返回最多5000条数据

```json
{
    "data": {
        "traceList": {
            "total": 17,
            "traces": [
                {
                    "spans": [
                        {
                            "operationName": "callQuery",
                            "startTime": 1534915192065927,
                            "tags": [
                                {
                                    "key": "sampler.type",
                                    "type": "string",
                                    "value": "const"
                                },
                                {
                                    "key": "sampler.param",
                                    "type": "bool",
                                    "value": true
                                },
                                {
                                    "key": "span.kind",
                                    "type": "string",
                                    "value": "client"
                                },
                                {
                                    "key": "http.url",
                                    "type": "string",
                                    "value": "http://localhost:12581/query/49"
                                },
                                {
                                    "key": "http.method",
                                    "type": "string",
                                    "value": "GET"
                                },
                                {
                                    "key": "fromId",
                                    "type": "int64",
                                    "value": 49
                                }
                            ]
                        },
                        {
                            "operationName": "forword-query",
                            "startTime": 1534915192067234,
                            "tags": [
                                {
                                    "key": "span.kind",
                                    "type": "string",
                                    "value": "server"
                                }
                            ]
                        },
                        {
                            "operationName": "do-query",
                            "startTime": 1534915192067356,
                            "tags": [
                                {
                                    "key": "span.kind",
                                    "type": "string",
                                    "value": "server"
                                }
                            ]
                        }
                    ],
                    "traceID": "4af90535eefe778c"
                }
            ]
        }
    }
}
```

# 查询某个trace的详细

```json
{
	"query": "query test($traceId: ID!) { trace(traceId: $traceId) { traceID spans { operationName tags { key type value } } } }",
	"variables": {"traceId":"52c45831057fdff4" }
}
```

返回结果跟trace列表相同


# 查询吞吐量

按一个时间范围，服务（Operation）查询这个时间范围内，按分钟统计的吞吐量（即查询span的个数）

```json
{
	"query": "query q($serviceId: ID!, $duration: Duration!) { throughput(serviceId: $serviceId, duration: $duration) { trendList }}",
	"variables": { "duration":{ "start": "2018-08-22 13:15", "end": "2018-08-22 13:30", "step": "MUNITE"  }, "serviceId":"/sample"  }
}
```

返回结果示例

```json
{"data":{"throughout":{"trendList":[0,0,0,0,4,7,4,7,4,8,7,7,7,2,0,0]}}}
```

# 查询平均响应时间

按一个时间范围，服务（Operation）查询这个时间范围内，按分钟统计的平均响应时间

```json
{
	"query": "query q($serviceId: ID!, $duration: Duration!) { responseTime(serviceId: $serviceId, duration: $duration) { trendList }}",
	"variables": { "duration":{ "start": "2018-08-22 13:15", "end": "2018-08-22 13:30", "step": "MUNITE"  }, "serviceId":"/sample"  }
}
```

返回结果示例

```json
{"data":{"responseTime":{"trendList":[0,0,0,0,1002,1297,1343,1336,1631,1663,1598,1433,1715,1190,0,0]}}}
```

# 查询数据库信息

查询当前数据库类型的服务数量和对端信息

```json
{
	"query": "query q($duration: Duration!) { dbInfo(duration: $duration) { count, peers } }",
	"variables": { "duration":{ "start": "2018-08-22 13:15", "end": "2018-08-22 13:30", "step": "MUNITE"  } }
}
```

返回结果示例

```json
{"data":{"dbInfo":{"count":2,"peers":["192.168.31.102:3306","192.168.31.103:3306"]}}}
```

# 查询缓存信息

查询当前缓存类型的服务数量和对端信息

```json
{
	"query": "query q($duration: Duration!) { cacheInfo(duration: $duration) { count, peers } }",
	"variables": { "duration":{ "start": "2018-08-22 13:15", "end": "2018-08-22 13:30", "step": "MUNITE"  } }
}
```

返回结果示例

```json
{"data":{"cacheInfo":{"count":2,"peers":["192.168.31.102:6397","192.168.31.100:6397"]}}}
```