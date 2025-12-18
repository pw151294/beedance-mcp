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
	)
}

func ListTraces(variable ListTracesVariable, headers map[string]string) (string, error) {
	// 1. 发送graphql请求
	graphqlResp, err := graphql.DoGraphqlRequest[ListTracesVariable, ListTracesResponse](listTracesGraphqlQuery, headers, variable)
	if err != nil {
		loggers.Error("send graphql request failed", zap.Any("variable", variable), zap.Error(err))
		return "", err
	}

	// 2. 将工具调用的结果转换成白话文
	tracesData := graphqlResp.Data
	loggers.Info("tracesData", zap.Any("tracesData", tracesData))
	return convert2Message(tracesData.Data), nil
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
	listTracesMessage, err := ListTraces(variable, headers)
	if err != nil {
		loggers.Error("list traces message failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultError("发送graphql请求失败：" + err.Error()), nil
	}
	loggers.Info("invoke list traces tool success", zap.Any("message", listTracesMessage))
	return mcp.NewToolResultText(listTracesMessage), nil
}
