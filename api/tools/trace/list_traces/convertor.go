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

func convert2Variable(request mcp.CallToolRequest) (ListTracesVariable, error) {
	workspaceId := request.Header.Get(tools.WorkspaceIdHeaderName)
	if workspaceId == "" {
		loggers.Error("parse workspaceId from request headers failed,", zap.Any("headers", request.Header))
		return ListTracesVariable{}, errors.New("获取工作空间ID参数失败")
	}

	// 1.服务名称
	serviceName, err := request.RequireString(tools.ServiceNameParamName)
	if err != nil {
		loggers.Error("parse service name from request failed", zap.Any("request", request), zap.Error(err))
		return ListTracesVariable{}, fmt.Errorf("获取服务名称参数失败：%v", err)
	}
	serviceIds := list_services.ConvertServiceNames2IDs(request, workspaceId, []string{serviceName})
	serviceId := serviceIds[0]
	endpointName := request.GetString(endpointNameParamName, "")
	endpointId := convertor.ConvertServiceIDAndEndpointName2EndpointID(serviceId, endpointName)

	// 2.开始时间 链路状态
	start := request.GetString(tools.StartParamName, "")
	duration, err := timeutils.BuildDuration(start)
	if err != nil {
		loggers.Error("parse duration from request failed", zap.Any("start", start), zap.Error(err))
		return ListTracesVariable{}, fmt.Errorf("获取查询时间参数失败：%v", err)
	}
	traceState := request.GetString(traceStateParamName, "ALL")

	// 3.查询条件
	condition := Condition{}
	condition.TraceState = traceState
	condition.ServiceId = serviceId
	condition.EndpointId = endpointId
	condition.QueryOrder = queryOrder
	condition.QueryDuration = duration
	condition.Paging = Paging{PageNum: pageNum, PageSize: pageSize}

	return ListTracesVariable{Condition: condition}, nil
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
