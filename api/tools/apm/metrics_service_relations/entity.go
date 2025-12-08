package metrics_service_relations

import "beedance-mcp/api/tools/apm"

// ServiceRelationMetricsVariables 查询服务关系客户端指标的变量
type ServiceRelationMetricsVariables struct {
	WorkspaceID string       `json:"-"`
	Duration    apm.Duration `json:"duration"`
	IDs         []string     `json:"ids"`
	M0          string       `json:"m0"` // service_relation_client_cpm
	M1          string       `json:"m1"` // service_relation_client_resp_time
}

// ServiceRelationClientMetricsResponse 查询服务关系指标的响应
type ServiceRelationClientMetricsResponse struct {
	ServiceRelationClientCPM      apm.MetricValues `json:"service_relation_client_cpm"`
	ServiceRelationClientRespTime apm.MetricValues `json:"service_relation_client_resp_time"`
}

// ServiceRelationServerMetricsResponse 查询服务关系服务端指标的响应
type ServiceRelationServerMetricsResponse struct {
	ServiceRelationServerRespTime apm.MetricValues `json:"service_relation_server_resp_time"`
	ServiceRelationServerCPM      apm.MetricValues `json:"service_relation_server_cpm"`
}
