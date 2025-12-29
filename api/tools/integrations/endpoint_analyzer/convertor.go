package endpoint_analyzer

import (
	"beedance-mcp/api/tools"
	"beedance-mcp/internal/pkg/convertor"
	"beedance-mcp/internal/pkg/extractor"
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

func convert2EndpointTracesMcpRequest(request mcp.CallToolRequest, traceState string) (map[string]mcp.CallToolRequest, error) {
	endpointIds, err := request.RequireStringSlice(endpointIdsParamName)
	if err != nil || len(endpointIds) == 0 {
		loggers.Error("parse endpointIds from mcp request failed", zap.Any("request", request), zap.Error(err))
		return nil, err
	}

	serviceName2EndpointIds := make(map[string][]string)
	for _, endpointId := range endpointIds {
		serviceId, _ := convertor.ConvertEndpointID2ServiceIDAndEndpointName(endpointId)
		serviceName := convertor.ConvertServiceID2Name(serviceId)
		if endpointIDs, ok := serviceName2EndpointIds[serviceName]; !ok {
			serviceName2EndpointIds[serviceName] = []string{endpointId}
		} else {
			serviceName2EndpointIds[serviceName] = append(endpointIDs, endpointId)
		}
	}

	serviceName2EndpointTracesMcpRequest := make(map[string]mcp.CallToolRequest)
	for serviceName, endpointIds := range serviceName2EndpointIds {
		mcpRequest := copyMcpRequest(request)
		arguments := make(map[string]any)
		arguments[tools.ServiceNameParamName] = serviceName
		arguments[endpointIdsParamName] = endpointIds
		arguments[tools.StartParamName] = request.GetString(tools.StartParamName, "")
		arguments[traceStateParamName] = traceState
		mcpRequest.Params.Arguments = arguments
		serviceName2EndpointTracesMcpRequest[serviceName] = mcpRequest
	}
	return serviceName2EndpointTracesMcpRequest, nil
}

func convert2DetailTracesMcpRequest(request mcp.CallToolRequest, endpointTraces []extractor.EndpointTrace) map[string]mcp.CallToolRequest {
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
