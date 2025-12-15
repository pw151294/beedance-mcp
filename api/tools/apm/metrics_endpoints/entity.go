package metrics_endpoints

import (
	"beedance-mcp/api/tools"
)

type MetricsEndpointVariable struct {
	Duration     tools.Duration `json:"duration"`
	Condition0   Condition      `json:"condition0"`
	GraphqlQuery string         `json:"-"`
}

type Condition struct {
	Name          string `json:"name"`
	ParentService string `json:"parentService"`
	Normal        bool   `json:"normal"`
	Scope         string `json:"scope"`
	TopN          int    `json:"topN"`
	Order         string `json:"order"`
}

type MetricValue struct {
	Name  string `json:"name"`
	Id    string `json:"id"`
	Value string `json:"value"`
}

type MetricsEndpointResponse struct {
	MetricsEndpointRespTime []MetricValue `json:"endpoint_resp_time0"`
	MetricsEndpointCpm      []MetricValue `json:"endpoint_cpm0"`
	MetricsEndpointSla      []MetricValue `json:"endpoint_sla0"`
}
