package list_traces

import "beedance-mcp/api/tools"

type ListEndpointsVariable struct {
	ServiceId string         `json:"serviceId"`
	Duration  tools.Duration `json:"duration"`
	Keyword   string         `json:"keyword"`
}
type ListEndpointResponse struct {
	Pods []Pod `json:"pods"`
}
type Pod struct {
	Id    string `json:"id"`
	Value string `json:"value"`
	Label string `json:"label"`
}

type Paging struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

type Condition struct {
	QueryDuration tools.Duration `json:"queryDuration"`
	TraceState    string         `json:"traceState"`
	QueryOrder    string         `json:"queryOrder"`
	Paging        Paging         `json:"paging"`
	ServiceId     string         `json:"serviceId"`
	EndpointId    string         `json:"endpointId,omitempty"`
}

type ListTracesVariable struct {
	Condition Condition `json:"condition"`
}

type Trace struct {
	Key           string   `json:"key" jsonschema_description:"链路片段的唯一标识符"`
	EndpointNames []string `json:"endpointNames" jsonschema_description:"链路中涉及的端点名称列表"`
	Duration      int      `json:"duration" jsonschema_description:"链路持续时间（毫秒）"`
	Start         string   `json:"start" jsonschema_description:"链路开始的时间戳"`
	IsError       bool     `json:"isError" jsonschema_description:"指示链路是否包含错误"`
	TraceIds      []string `json:"traceIds" jsonschema_description:"与此片段关联的链路ID列表"`
}

type TracesData struct {
	Traces []Trace `json:"traces" jsonschema_description:"找到的链路片段列表"`
}

type ListTracesResponse struct {
	Data TracesData `json:"data" jsonschema_description:"链路查询结果的容器"`
}
