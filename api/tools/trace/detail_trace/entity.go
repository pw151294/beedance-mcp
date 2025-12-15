package detail_trace

type DetailTraceVariable struct {
	TraceId string `json:"traceId"`
}

type Span struct {
	TraceId             string        `json:"traceId"`
	SegmentId           string        `json:"segmentId"`
	SpanId              int           `json:"spanId"`
	ParentSpanId        int           `json:"parentSpanId"`
	Refs                []interface{} `json:"refs"`
	ServiceCode         string        `json:"serviceCode"`
	ServiceInstanceName string        `json:"serviceInstanceName"`
	StartTime           int64         `json:"startTime"`
	EndTime             int64         `json:"endTime"`
	EndpointName        string        `json:"endpointName"`
	Type                string        `json:"type"`
	Peer                string        `json:"peer"`
	Component           string        `json:"component"`
	IsError             bool          `json:"isError"`
	Layer               string        `json:"layer"`
	Tags                []Tag         `json:"tags"`
	Logs                []Log         `json:"logs"`
}

type TraceDetail struct {
	Spans []Span `json:"spans"`
}

type Log struct {
	Time int64
	Data []LogTag
}
type LogTag struct {
	Key   string
	Value string
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type DetailTraceResponse struct {
	TraceDetail TraceDetail `json:"trace"`
}
