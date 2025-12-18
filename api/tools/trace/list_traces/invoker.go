package list_traces

import (
	"beedance-mcp/api/tools"
	"beedance-mcp/internal/pkg/graphql"
	"beedance-mcp/pkg/httputils"
	"beedance-mcp/pkg/loggers"
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

func ListTracesToolSchema() mcp.Tool {
	return mcp.NewTool(
		listTracesToolName,
		mcp.WithDescription(listTracesToolDesc),
		mcp.WithString(tools.StartParamName, mcp.Description(tools.StartParamDesc)),
		mcp.WithString(traceStateParamName, mcp.Description(traceStateParamDesc)),
		mcp.WithString(tools.ServiceNameParamName, mcp.Required(), mcp.Description(tools.ServiceNameParamDesc)),
		mcp.WithString(endpointNameParamName, mcp.Description(endpointNameParamDesc)),
		mcp.WithOutputSchema[ListTracesResponse](),
	)
}

func ListTraces(variable ListTracesVariable, headers map[string]string) (ListTracesResponse, error) {
	// 1. 发送graphql请求
	graphqlResp, err := graphql.DoGraphqlRequest[ListTracesVariable, ListTracesResponse](listTracesGraphqlQuery, headers, variable)
	if err != nil {
		loggers.Error("send graphql request failed", zap.Any("variable", variable), zap.Error(err))
		return ListTracesResponse{}, err
	}
	return graphqlResp.Data, nil
}

func InvokeListTracesTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 1. 转换请求参数
	variable, err := convert2Variable(request)
	if err != nil {
		loggers.Error("convert request to variable failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultError("获取graphql请求参数失败：" + err.Error()), nil
	}

	// 2. 发送graphql请求
	headers, err := httputils.BuildHeaders(request)
	if err != nil {
		loggers.Error("build headers from request failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultError("获取graphql请求头失败：" + err.Error()), nil
	}
	listTracesResp, err := ListTraces(variable, headers)
	if err != nil {
		loggers.Error("invoke list traces failed", zap.Error(err))
		return mcp.NewToolResultError("调用list_taces工具失败：" + err.Error()), nil
	}

	// 3. 使用结构化输出
	loggers.Info("list traces response", zap.Any("list traces response", listTracesResp))
	message := convert2Message(listTracesResp.Data)
	loggers.Info("tool invoke success", zap.Any("message", message))
	return mcp.NewToolResultStructured(listTracesResp, message), nil
}
