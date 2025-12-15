package metrics_endpoints

import (
	"beedance-mcp/api/tools"
	"beedance-mcp/internal/pkg/graphql"
	"beedance-mcp/pkg/httputils"
	"beedance-mcp/pkg/loggers"
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

func MetricsEndpointsToolSchema() mcp.Tool {
	return mcp.NewTool(
		metricsEndpointsToolName,
		mcp.WithDescription(metricsEndpointsToolDesc),
		mcp.WithString(tools.ServiceNameParamName, mcp.Required(), mcp.Description(tools.ServiceNameParamDesc)),
		mcp.WithString(metricsNameParamName, mcp.Required(), mcp.Description(metricsNameParamDesc)),
		mcp.WithString(tools.StartParamName, mcp.Description(tools.StartParamDesc)),
		mcp.WithNumber(topNParamName, mcp.Description(topNParamDesc)),
	)
}

func InvokeMetricsEndpointsTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 1. 转换请求参数
	variable, err := convert2Variable(request)
	if err != nil {
		loggers.Error("convert request to graphql variable failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultError("构建graphql请求变量失败：" + err.Error()), nil
	}

	// 2. 发送graphql请求
	headers, err := httputils.BuildHeaders(request)
	if err != nil {
		loggers.Error("build headers failed", zap.Error(err))
		return mcp.NewToolResultError("构建请求头失败：" + err.Error()), nil
	}
	graphqlResp, err := graphql.DoGraphqlRequest[MetricsEndpointVariable, MetricsEndpointResponse](variable.GraphqlQuery, headers, variable)
	if err != nil {
		loggers.Error("execute graphql request failed", zap.Error(err))
		return mcp.NewToolResultError("发送graphql请求失败：" + err.Error()), nil
	}

	// 3. 将查询到的结果转换为白话文
	metricsResp := graphqlResp.Data
	loggers.Info("send graphql request success", zap.Any("metricsResp", metricsResp))
	message := convert2Message(metricsResp)
	loggers.Info("invoke metrics endpoints tool success", zap.String("tool invoke message", message))
	return mcp.NewToolResultText(message), nil
}
