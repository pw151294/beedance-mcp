package tools

// Duration 时间范围参数
type Duration struct {
	Start string `json:"start"` // 格式: "2025-12-05 0934"
	End   string `json:"end"`   // 格式: "2025-12-05 1004"
	Step  string `json:"step"`  // 如: "MINUTE"
}

const WorkspaceIdHeaderName = "Workspaceid"
const TokenHeaderName = "Token"

const (
	ServiceNamesParamName = "serviceNames"
	StartParamName        = "start"
	StartParamDesc        = "查询时间范围的起始时间，遵循'YYYY-MM-DD HH::mm:ss'的格式"
	ServiceNamesParamDesc = "服务名称列表"
)
