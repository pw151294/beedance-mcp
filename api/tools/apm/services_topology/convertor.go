package services_topology

import (
	"beedance-mcp/api/tools/apm"
	"beedance-mcp/api/tools/apm/list_services"
	"beedance-mcp/pkg/loggers"
	"beedance-mcp/pkg/timeutils"
	"bytes"
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

func convert2TopoVariables(request mcp.CallToolRequest) (ServiceTopologyVariables, error) {
	workspaceId := request.Header.Get(apm.WorkspaceIdHeaderName)
	if workspaceId == "" {
		loggers.Error("parse workspaceId from header failed", zap.Any("headers", request.Header))
		return ServiceTopologyVariables{}, errors.New("请求头未携带工作空间ID")
	}
	serviceNames, err := request.RequireStringSlice(apm.ServiceNamesParamName)
	if err != nil {
		loggers.Error("parse serviceNames failed", zap.Error(err))
		return ServiceTopologyVariables{}, fmt.Errorf("服务名称列表参数错误：%w", err)
	}
	start := request.GetString(apm.StartParamName, "")
	duration, err := timeutils.BuildDuration(start)
	if err != nil {
		loggers.Error("build duration failed", zap.String("start", start), zap.Error(err))
		return ServiceTopologyVariables{}, fmt.Errorf("构建duration参数错误：%w", err)
	}

	variables := ServiceTopologyVariables{}
	variables.WorkspaceID = workspaceId
	variables.Duration = duration
	variables.IDs = list_services.ServiceIDs(workspaceId, serviceNames)
	return variables, nil
}

func convert2Message(workspaceId string, serviceTopoResp ServiceTopologyResponse) string {
	topology := serviceTopoResp.Topology
	nodes := topology.Nodes
	calls := topology.Calls

	var toolInvokeMessageBuffer bytes.Buffer
	if len(nodes) == 0 || len(calls) == 0 {
		toolInvokeMessageBuffer.WriteString("为查询到服务调用关系，服务调用拓扑为空")
	} else {
		toolInvokeMessageBuffer.WriteString("服务调用拓扑如下图所示：\n")
		for _, call := range calls {
			srcNode := GetNode(workspaceId, call.Source)
			targetNode := GetNode(workspaceId, call.Target)
			srcName := list_services.ServiceName(workspaceId, call.Source)
			targetName := list_services.ServiceName(workspaceId, call.Target)
			toolInvokeMessageBuffer.WriteString(fmt.Sprintf(serviceCallInfoPattern, srcName, srcNode.Type, targetName, targetNode.Type))
		}
	}

	return toolInvokeMessageBuffer.String()
}
