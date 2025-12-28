package service_analyzer

import (
	"beedance-mcp/api/tools"
	"beedance-mcp/pkg/loggers"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

func copyMcpRequest(request mcp.CallToolRequest) mcp.CallToolRequest {
	newRequest := mcp.CallToolRequest{
		Request: request.Request,
		Header:  make(map[string][]string),
		Params: mcp.CallToolParams{
			Name:      request.Params.Name,
			Arguments: make(map[string]interface{}),
		},
	}

	// 深拷贝 Header
	for k, v := range request.Header {
		headerCopy := make([]string, len(v))
		copy(headerCopy, v)
		newRequest.Header[k] = headerCopy
	}

	return newRequest
}

func convert2MetricsEndpointsMcpRequest(request mcp.CallToolRequest, metricsName string) (mcp.CallToolRequest, error) {
	mcpRequest := copyMcpRequest(request)

	serviceName, err := request.RequireString(tools.ServiceNameParamName)
	if err != nil {
		loggers.Error("parse serviceName from mcp call tool request failed", zap.Any("request", request), zap.Error(err))
		return mcpRequest, err
	}

	arguments := make(map[string]any)
	arguments[tools.ServiceNameParamName] = serviceName
	arguments[tools.StartParamName] = request.GetString(tools.StartParamName, "")
	arguments[metricsNameParamName] = metricsName
	arguments[topNParamName] = 5
	mcpRequest.Params.Arguments = arguments
	return mcpRequest, nil
}

func convert2EndpointTracesMcpRequest(request mcp.CallToolRequest, endpointIds []string, metricsName string) mcp.CallToolRequest {
	mcpRequest := copyMcpRequest(request)

	serviceName, _ := request.RequireString(tools.ServiceNameParamName)
	arguments := make(map[string]any)
	arguments[tools.ServiceNameParamName] = serviceName
	arguments[tools.StartParamName] = request.GetString(tools.StartParamName, "")
	arguments[endpointIdsParamName] = endpointIds
	if metricsName == endpointSlaMetricsName {
		arguments[traceStateParamName] = "ERROR"
	}
	mcpRequest.Params.Arguments = arguments

	return mcpRequest
}

func convert2DetailTracesMcpRequest(request mcp.CallToolRequest, endpointTraces []EndpointTrace) map[string]mcp.CallToolRequest {
	endpointName2McpRequest := make(map[string]mcp.CallToolRequest)
	for _, endpointTrace := range endpointTraces {
		mcpRequest := copyMcpRequest(request)
		arguments := make(map[string]any)
		traceIds := endpointTrace.TraceIds
		arguments[traceIDsParamName] = traceIds[:min(3, len(traceIds))] // 这里只选前3个有代表性的链路来分析
		mcpRequest.Params.Arguments = arguments
		endpointName2McpRequest[endpointTrace.EndpointName] = mcpRequest
	}
	return endpointName2McpRequest
}
