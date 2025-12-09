package services_topology

import (
	"beedance-mcp/api/tools/apm"
	"beedance-mcp/api/tools/apm/list_services"
	"beedance-mcp/internal/pkg/convertor"
	"beedance-mcp/pkg/loggers"
	"beedance-mcp/pkg/table"
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
	variables.IDs = list_services.ConvertServiceNames2IDs(request, workspaceId, serviceNames)
	return variables, nil
}

func convert2Tables(workspaceId string, serviceTopoResp ServiceTopologyResponse) (nodeRegister *table.Table[string, string, Node], callRegister *table.Table[string, string, Call]) {
	nodeRegister = table.NewTable[string, string, Node]()
	callRegister = table.NewTable[string, string, Call]()
	topology := serviceTopoResp.Topology
	nodes := topology.Nodes
	calls := topology.Calls

	if len(nodes) > 0 {
		for _, node := range nodes {
			nodeRegister.Put(workspaceId, node.ID, node)
		}
	}

	if len(calls) > 0 {
		for _, call := range calls {
			callRegister.Put(workspaceId, call.ID, call)
		}
	}

	return
}

func convert2Message(workspaceId string, serviceTopoResp ServiceTopologyResponse) string {
	nodeRegister, callRegister := convert2Tables(workspaceId, serviceTopoResp)

	var toolInvokeMessageBuffer bytes.Buffer
	if nodeRegister.Size() == 0 || callRegister.Size() == 0 {
		toolInvokeMessageBuffer.WriteString("未查询到服务调用关系，服务调用拓扑为空")
	} else {
		toolInvokeMessageBuffer.WriteString("服务调用拓扑如下：\n")
		for _, call := range serviceTopoResp.Topology.Calls {
			srcNode, _ := nodeRegister.Get(workspaceId, call.Source)
			tgtNode, _ := nodeRegister.Get(workspaceId, call.Target)
			srcName := convertor.ConvertID2Name(call.Source)
			tgtName := convertor.ConvertID2Name(call.Target)
			toolInvokeMessageBuffer.WriteString(fmt.Sprintf(serviceCallInfoPattern, srcName, srcNode.Type, tgtName, tgtNode.Type))
		}
	}

	return toolInvokeMessageBuffer.String()
}
