package services_topology

import (
	"beedance-mcp/api/tools"
)

// Node 拓扑节点
type Node struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	IsReal bool   `json:"isReal"`
}

// Call 调用关系
type Call struct {
	ID           string   `json:"id"`
	Source       string   `json:"source"`
	DetectPoints []string `json:"detectPoints"`
	Target       string   `json:"target"`
}

// Topology 拓扑数据
type Topology struct {
	Nodes []Node `json:"nodes"`
	Calls []Call `json:"calls"`
}

type ServiceTopologyVariables struct {
	WorkspaceID string         `json:"-"`
	IDs         []string       `json:"serviceIds"`
	Duration    tools.Duration `json:"duration"`
}

type ServiceTopologyResponse struct {
	Topology Topology `json:"topology"`
}
