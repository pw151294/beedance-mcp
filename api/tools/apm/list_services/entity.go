package list_services

// Service 服务信息
type Service struct {
	ID        string   `json:"id"`
	Value     string   `json:"value"` // 服务名称
	Label     string   `json:"label"` // 短名称
	Group     string   `json:"group"`
	Layers    []string `json:"layers"`
	Normal    bool     `json:"normal"`
	GroupName string   `json:"groupName"`
}

// ListServicesVariables 查询服务列表的变量
type ListServicesVariables struct {
	Layer       string `json:"layer"`
	WorkspaceID string `json:"workspaceId,omitempty"`
}

// ListServicesResponse 查询服务列表的响应
type ListServicesResponse struct {
	Services []Service `json:"services"`
}
