package detail_trace

type DetailTraceVariable struct {
	TraceId string `json:"traceId"`
}

type Span struct {
	TraceId             string        `json:"traceId" jsonschema_description:"链路ID"`
	SegmentId           string        `json:"segmentId" jsonschema_description:"片段ID"`
	SpanId              int           `json:"spanId" jsonschema_description:"跨度ID"`
	ParentSpanId        int           `json:"parentSpanId" jsonschema_description:"父跨度ID"`
	Refs                []interface{} `json:"refs" jsonschema_description:"引用列表"`
	ServiceCode         string        `json:"serviceCode" jsonschema_description:"服务代码"`
	ServiceInstanceName string        `json:"serviceInstanceName" jsonschema_description:"服务实例名称"`
	StartTime           int64         `json:"startTime" jsonschema_description:"开始时间（毫秒）"`
	EndTime             int64         `json:"endTime" jsonschema_description:"结束时间（毫秒）"`
	EndpointName        string        `json:"endpointName" jsonschema_description:"接口名称"`
	Type                string        `json:"type" jsonschema_description:"跨度类型"`
	Peer                string        `json:"peer" jsonschema_description:"对端地址"`
	Component           string        `json:"component" jsonschema_description:"组件名称"`
	IsError             bool          `json:"isError" jsonschema_description:"是否包含错误"`
	Layer               string        `json:"layer" jsonschema_description:"层级"`
	Tags                []Tag         `json:"tags" jsonschema_description:"标签列表"`
	Logs                []Log         `json:"logs" jsonschema_description:"日志列表"`
}

type TraceDetail struct {
	Spans []Span `json:"spans" jsonschema_description:"链路跨度列表"`
}

type Log struct {
	Time int64    `json:"time" jsonschema_description:"日志时间戳"`
	Data []LogTag `json:"data" jsonschema_description:"日志数据键值对列表"`
}
type LogTag struct {
	Key   string `json:"key" jsonschema_description:"键"`
	Value string `json:"value" jsonschema_description:"值"`
}

type Tag struct {
	Key   string `json:"key" jsonschema_description:"标签键"`
	Value string `json:"value" jsonschema_description:"标签值"`
}

type DetailTraceResponse struct {
	TraceDetail TraceDetail `json:"trace" jsonschema_description:"链路详情数据"`
}
