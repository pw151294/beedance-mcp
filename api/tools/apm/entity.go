package apm

// Duration 时间范围参数
type Duration struct {
	Start string `json:"start"` // 格式: "2025-12-05 0934"
	End   string `json:"end"`   // 格式: "2025-12-05 1004"
	Step  string `json:"step"`  // 如: "MINUTE"
}

// MetricValue 指标值
type MetricValue struct {
	ID    string `json:"id"`
	Value int64  `json:"value"`
}

// MetricValues 指标值列表
type MetricValues struct {
	Values []MetricValue `json:"values"`
}

const (
	WorkspaceIdParamName = "workspaceId"
	WorkspaceIdParamDesc = "工作空间ID"
	TokenParamName       = "token"
	TokenParamDesc       = "用户身份信息校验token"
)
