package trace

import "beedance-mcp/api/tools"

type EndpointQueryRequest struct {
	ServiceId string         `json:"serviceId"`
	Duration  tools.Duration `json:"duration"`
}

type TraceQueryRequest struct {
	TraceId string `json:"traceId"`
}

type Pod struct {
	Id    string `json:"id"`
	Value string `json:"value"`
	Label string `json:"label"`
}

type EndpointQueryRespData struct {
	Pods []Pod `json:"pods"`
}

type EndpointQueryResp struct {
	Data EndpointQueryRespData `json:"data"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
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
	Logs                []interface{} `json:"logs"`
	AttachedEvents      []interface{} `json:"attachedEvents"`
}

type TraceDetail struct {
	Spans []Span `json:"spans"`
}

type TraceQueryRespData struct {
	Trace TraceDetail `json:"trace"`
}

type TraceQueryResp struct {
	Data TraceQueryRespData `json:"data"`
}
