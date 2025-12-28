package list_traces

import (
	"beedance-mcp/api/tools"
	"beedance-mcp/api/tools/apm/list_services"
	"beedance-mcp/internal/pkg/convertor"
	"beedance-mcp/pkg/loggers"
	"beedance-mcp/pkg/timeutils"
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

func convert2ListEndpointsVariables(request mcp.CallToolRequest) ([]ListEndpointsVariable, error) {
	workspaceId := request.Header.Get(tools.WorkspaceIdHeaderName)
	if workspaceId == "" {
		loggers.Error("parse workspaceId from request headers failed,", zap.Any("headers", request.Header))
		return nil, errors.New("获取工作空间ID参数失败")
	}

	// 1. 获取服务名称列表
	serviceNames := request.GetStringSlice(tools.ServiceNamesParamName, make([]string, 0, 0))
	var serviceIds []string
	if len(serviceNames) > 0 {
		serviceIds = list_services.ConvertServiceNames2IDs(request, workspaceId, serviceNames)
	} else {
		listServicesResp, err := list_services.ListServices(request)
		if err != nil {
			loggers.Error("list services failed,", zap.Any("request", request), zap.Error(err))
			return nil, fmt.Errorf("获取服务列表失败：%w", err)
		}
		services := listServicesResp.Services
		for _, svc := range services {
			serviceIds = append(serviceIds, svc.ID)
		}
	}

	// 2. 获取开始时间
	start := request.GetString(tools.StartParamName, "")
	duration, err := timeutils.BuildDuration(start)
	if err != nil {
		loggers.Error("parse duration from request failed", zap.Any("start", start), zap.Error(err))
		return nil, fmt.Errorf("获取查询时间参数失败：%v", err)
	}

	// 3. 构建查询条件
	variables := make([]ListEndpointsVariable, 0, len(serviceIds))
	for _, svcId := range serviceIds {
		variable := ListEndpointsVariable{}
		variable.ServiceId = svcId
		variable.Duration = duration
		variables = append(variables, variable)
	}
	return variables, nil
}

func convert2ListTracesVariables(request mcp.CallToolRequest) ([]ListTracesVariable, error) {
	workspaceId := request.Header.Get(tools.WorkspaceIdHeaderName)
	if workspaceId == "" {
		loggers.Error("parse workspaceId from request headers failed,", zap.Any("headers", request.Header))
		return nil, errors.New("获取工作空间ID参数失败")
	}

	serviceName, err := request.RequireString(tools.ServiceNameParamName)
	if err != nil {
		loggers.Error("parse serviceName from request failed,", zap.Any("request", request), zap.Error(err))
		return nil, fmt.Errorf("获取服务名称参数失败：%w", err)
	}
	serviceIds := list_services.ConvertServiceNames2IDs(request, workspaceId, []string{serviceName})
	serviceId := serviceIds[0]

	endpointIds, err := request.RequireStringSlice(endpointIdsParamName)
	if err != nil || len(endpointIds) == 0 {
		loggers.Error("parse endpointIds from request failed,", zap.Any("request", request), zap.Error(err))
		return nil, fmt.Errorf("获取接口ID列表参数失败：%w", err)
	}
	start := request.GetString(tools.StartParamName, "")
	duration, err := timeutils.BuildDuration(start)
	if err != nil {
		loggers.Error("parse duration from request failed", zap.Any("start", start), zap.Error(err))
		return nil, fmt.Errorf("获取查询时间参数失败：%v", err)
	}
	traceState := request.GetString(traceStateParamName, "ALL")

	variables := make([]ListTracesVariable, 0, len(endpointIds))
	for _, endpointId := range endpointIds {
		condition := Condition{}
		condition.TraceState = traceState
		condition.ServiceId = serviceId
		condition.EndpointId = endpointId
		condition.QueryOrder = queryOrder
		condition.QueryDuration = duration
		condition.Paging = Paging{PageNum: pageNum, PageSize: pageSize}

		variable := ListTracesVariable{Condition: condition}
		variables = append(variables, variable)
	}
	return variables, nil
}

func convertPod2Message(pod Pod) string {
	endpointId := pod.Id
	serviceId, endpointName := convertor.ConvertEndpointID2ServiceIDAndEndpointName(endpointId)
	serviceName := convertor.ConvertServiceID2Name(serviceId)
	return fmt.Sprintf(endpointInfoPattern, endpointName, endpointId, serviceName, serviceId)
}

func convertEndpoints2Message(endPoints []string) string {
	if len(endPoints) == 0 {
		return "未发现任何接口"
	}

	var endpointMessageBuffer bytes.Buffer
	endpointMessageBuffer.WriteString("[")
	endpointMessageBuffer.WriteString(strings.Join(endPoints, "；"))
	endpointMessageBuffer.WriteString("]")
	return endpointMessageBuffer.String()
}
func convertTraceIds2Message(traceIds []string) string {
	if len(traceIds) == 0 {
		return "链路ID为空"
	}

	var traceIdMessageBuffer bytes.Buffer
	traceIdMessageBuffer.WriteString("[")
	traceIdMessageBuffer.WriteString(strings.Join(traceIds, ";"))
	traceIdMessageBuffer.WriteString("]")
	return traceIdMessageBuffer.String()
}

func convert2Message(tracesData TracesData) string {
	var toolInvokeMessageBuffer bytes.Buffer

	traces := tracesData.Traces
	if len(traces) == 0 {
		toolInvokeMessageBuffer.WriteString("未查询到服务的任何链路信息")
	} else {
		for _, trace := range traces {
			traceIdMessage := convertTraceIds2Message(trace.TraceIds)
			endpointsMessage := convertEndpoints2Message(trace.EndpointNames)
			traceState := convertor.ConvertBool2Desc(trace.IsError)
			toolInvokeMessageBuffer.WriteString(fmt.Sprintf(traceInfoPattern, traceIdMessage, endpointsMessage, trace.Duration, traceState))
		}
	}

	return toolInvokeMessageBuffer.String()
}
