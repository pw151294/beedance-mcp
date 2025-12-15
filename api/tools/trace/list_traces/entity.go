package list_traces

import "beedance-mcp/api/tools"

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
	Key           string   `json:"key"`
	EndpointNames []string `json:"endpointNames"`
	Duration      int      `json:"duration"`
	Start         string   `json:"start"`
	IsError       bool     `json:"isError"`
	TraceIds      []string `json:"traceIds"`
}

type TracesData struct {
	Traces []Trace `json:"traces"`
}

type ListTracesResponse struct {
	Data TracesData `json:"data"`
}
