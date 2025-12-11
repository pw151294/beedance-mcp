package metrics_service_relations

import (
	"beedance-mcp/api/tools/apm"
	"beedance-mcp/api/tools/apm/services_topology"
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

func convert2Variables(request mcp.CallToolRequest) (ServiceRelationMetricsVariables, error) {
	workspaceId := request.Header.Get(apm.WorkspaceIdHeaderName)
	if workspaceId == "" {
		loggers.Error("parse workspaceId from header failed", zap.Any("headers", request.Header))
		return ServiceRelationMetricsVariables{}, errors.New("请求头未携带工作空间ID")
	}
	serviceNames, err := request.RequireStringSlice(apm.ServiceNamesParamName)
	if err != nil {
		loggers.Error("parse serviceNames failed", zap.Error(err))
		return ServiceRelationMetricsVariables{}, fmt.Errorf("服务名称列表参数错误：%w", err)
	}
	start := request.GetString(apm.StartParamName, "")
	duration, err := timeutils.BuildDuration(start)
	if err != nil {
		loggers.Error("build duration failed", zap.String("start", start), zap.Error(err))
		return ServiceRelationMetricsVariables{}, fmt.Errorf("构建duration参数错误：%w", err)
	}

	variables := ServiceRelationMetricsVariables{}
	variables.WorkspaceID = workspaceId
	variables.Duration = duration
	variables.IDs = services_topology.ConvertServiceNames2CallIDs(request, workspaceId, serviceNames)
	return variables, nil
}

func convert2ClientVariables(request mcp.CallToolRequest) (ServiceRelationMetricsVariables, error) {
	variables, err := convert2Variables(request)
	if err != nil {
		return variables, fmt.Errorf("构建service_relation_client指标查询参数失败：%w", err)
	}
	variables.M0 = metricsClientM0Name
	variables.M1 = metricsClientM1Name
	return variables, nil
}

func convert2ServerVariables(request mcp.CallToolRequest) (ServiceRelationMetricsVariables, error) {
	variables, err := convert2Variables(request)
	if err != nil {
		return variables, fmt.Errorf("构建service_relation_server指标查询参数失败：%w", err)
	}
	variables.M0 = metricsServerM0Name
	variables.M1 = metricsServerM1Name
	return variables, nil
}

func convert2Table(clientResp ServiceRelationClientMetricsResponse, serverResp ServiceRelationServerMetricsResponse) *table.Table[string, string, int64] {
	metricsRegister := table.NewTable[string, string, int64]()

	clientCpms := clientResp.ServiceRelationClientCPM.Values
	for _, cpm := range clientCpms {
		metricsRegister.Put(cpm.ID, metricsClientM0Name, cpm.Value)
	}
	clientRts := clientResp.ServiceRelationClientRespTime.Values
	for _, rt := range clientRts {
		metricsRegister.Put(rt.ID, metricsClientM1Name, rt.Value)
	}

	serverCpms := serverResp.ServiceRelationServerCPM.Values
	for _, cpm := range serverCpms {
		metricsRegister.Put(cpm.ID, metricsServerM0Name, cpm.Value)
	}
	serverRts := serverResp.ServiceRelationServerRespTime.Values
	for _, rt := range serverRts {
		metricsRegister.Put(rt.ID, metricsServerM1Name, rt.Value)
	}

	return metricsRegister
}

func convert2Message(request mcp.CallToolRequest, clientResp ServiceRelationClientMetricsResponse, serverResp ServiceRelationServerMetricsResponse) string {
	workspaceId, _ := request.RequireString(apm.WorkspaceIdHeaderName)
	metricsRegister := convert2Table(clientResp, serverResp)
	id2Node := services_topology.CollectId2Node(request, workspaceId)

	var toolInvokeMessageBuffer bytes.Buffer
	if metricsRegister.Size() == 0 {
		toolInvokeMessageBuffer.WriteString("未查询到任何服务调用指标数据")
	} else {
		toolInvokeMessageBuffer.WriteString("服务调用指标数据如下：\n")
		callIds := metricsRegister.Rows()
		for _, callId := range callIds {
			srcId, tgtId := convertor.ConvertCallID2ServiceIDs(callId)
			srcNode, tgtNode := id2Node[srcId], id2Node[tgtId]
			srcName, tgtName := convertor.ConvertID2Name(srcId), convertor.ConvertID2Name(tgtId)

			clientCpm, _ := metricsRegister.Get(callId, metricsClientM0Name)
			serverRt, _ := metricsRegister.Get(callId, metricsServerM1Name)
			// todo 需要明确client/server数据的含义
			toolInvokeMessageBuffer.WriteString(fmt.Sprintf(serviceRelationInfoPattern, srcName, tgtName, srcNode.Type, tgtNode.Type, clientCpm, serverRt))
		}
	}

	return toolInvokeMessageBuffer.String()
}
