package services_topology

import (
	"beedance-mcp/internal/pkg/graphql"
	"beedance-mcp/pkg/httputils"
	"beedance-mcp/pkg/loggers"
	"beedance-mcp/pkg/table"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

var topoRegister *TopoRegister

type TopoRegister struct {
	id2Node *table.Table[string, string, Node] // workspaceId -> serviceID -> Node
	id2Call *table.Table[string, string, Call] // workspaceId -> callID -> Call
}

func newTopoRegister() *TopoRegister {
	return &TopoRegister{
		id2Node: table.NewTable[string, string, Node](),
		id2Call: table.NewTable[string, string, Call](),
	}
}
func ClearTopoRegister() {
	topoRegister = nil
}
func InitTopoRegister(request mcp.CallToolRequest) {
	if topoRegister != nil {
		return
	}
	topoRegister = newTopoRegister()
	variables, err := convert2TopoVariables(request)
	if err != nil {
		loggers.Error("convert to topo graph request variables failed", zap.Error(err), zap.Any("request", request))
		return
	}
	headers, err := httputils.BuildHeaders(request)
	if err != nil {
		loggers.Error("build graphql request headers failed", zap.Error(err), zap.Any("request", request))
		return
	}

	graphqlResp, err := graphql.DoGraphqlRequest[ServiceTopologyVariables, ServiceTopologyResponse](graphqlQuery, headers, variables)
	if err != nil {
		loggers.Error("send graphql request failed", zap.Error(err), zap.Any("request", request))
		return
	}
	topoResp := graphqlResp.Data
	loggers.Info("topo graph response: ", zap.Any("topoResp", topoResp))
	topoRegister.refresh(variables.WorkspaceID, topoResp)
}
func GetNode(workspaceId, svcID string) Node {
	return topoRegister.getNode(workspaceId, svcID)
}
func (tr *TopoRegister) refresh(workspaceId string, topoResp ServiceTopologyResponse) {
	topology := topoResp.Topology
	nodes := topology.Nodes
	calls := topology.Calls
	if len(nodes) > 0 {
		for _, node := range nodes {
			tr.id2Node.Put(workspaceId, node.ID, node)
		}
	}
	if len(calls) > 0 {
		for _, call := range calls {
			tr.id2Call.Put(workspaceId, call.ID, call)
		}
	}
}
func (tr *TopoRegister) getNode(workspaceId string, svcID string) Node {
	node, ok := tr.id2Node.Get(workspaceId, svcID)
	if !ok || node.ID == "" {
		loggers.Warn("get node from topo register failed",
			zap.String("workspaceId", workspaceId), zap.String("svcID", svcID), zap.Any("id2Node", tr.id2Node))
	}
	return node
}
