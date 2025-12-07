package metrics_services

import (
	"beedance-mcp/api/tools/apm"
	"beedance-mcp/api/tools/apm/list_services"
	"beedance-mcp/internal/pkg/graphql"
	"beedance-mcp/pkg/httputils"
	"beedance-mcp/pkg/loggers"
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

func MetricsServiceToolSchema() mcp.Tool {
	return mcp.NewTool(
		metricsServiceToolName,
		mcp.WithDescription(metricsServiceToolDescription),
		mcp.WithString(apm.WorkspaceIdParamName, mcp.Required(), mcp.Description(apm.WorkspaceIdParamDesc)),
		mcp.WithString(apm.TokenParamName, mcp.Required(), mcp.Description(apm.TokenParamDesc)),
		mcp.WithString(startParamName, mcp.Description(startParamDesc)),
		mcp.WithArray(serviceNamesParamName, mcp.Required(), mcp.Description(serviceNamesParamDesc)),
	)
}

func InvokeMetricsServicesTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 1. 转换请求参数
	list_services.InitServicesRegister(request)
	variables, err := convert2Variables(request)
	if err != nil {
		loggers.Error("convert request to graphql variables failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultError("获取graphql请求参数失败：" + err.Error()), nil
	}

	// 2. 发送graphql请求
	headers, err := httputils.BuildHeaders(request)
	if err != nil {
		loggers.Error("build headers failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultError("构建graphql请求头失败：" + err.Error()), nil
	}

	graphqlResp, err := graphql.DoGraphqlRequest[ServiceMetricsVariables, ServiceMetricsResponse](graphqlQuery, headers, variables)
	if err != nil {
		loggers.Error("call graphql request failed", zap.Any("variables", variables), zap.Any("headers", headers), zap.Error(err))
		return mcp.NewToolResultError("调用graphql接口失败：" + err.Error()), nil
	}

	// 3. 将工具调用结果转换成白话文
	loggers.Info("call graphql request success", zap.Any("metrics services", graphqlResp))
	message := convert2Message(variables.WorkspaceID, graphqlResp.Data)
	loggers.Info("tool invoke success", zap.String("message", message))
	return mcp.NewToolResultText(message), nil
}
