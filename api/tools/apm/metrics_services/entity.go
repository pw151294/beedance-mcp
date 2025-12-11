package metrics_services

import (
	"beedance-mcp/api/tools"
	"beedance-mcp/api/tools/apm"
)

// ServiceMetricsVariables 查询服务指标的变量
type ServiceMetricsVariables struct {
	WorkspaceID string         `json:"-"`
	Duration    tools.Duration `json:"duration"`
	IDs         []string       `json:"ids"`
	M0          string         `json:"m0"` // service_cpm
	M1          string         `json:"m1"` // service_sla
	M2          string         `json:"m2"` // service_resp_time
}

// ServiceMetricsRequest 查询服务指标的请求
type ServiceMetricsRequest struct {
	Query     string                  `json:"query"`
	Variables ServiceMetricsVariables `json:"variables"`
}

// ServiceMetricsResponse 查询服务指标的响应
type ServiceMetricsResponse struct {
	ServiceCPM      apm.MetricValues `json:"service_cpm"`
	ServiceSLA      apm.MetricValues `json:"service_sla"`
	ServiceRespTime apm.MetricValues `json:"service_resp_time"`
}
