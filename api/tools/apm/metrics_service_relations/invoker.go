package metrics_service_relations

import (
	"beedance-mcp/api/tools"
	"beedance-mcp/internal/pkg/graphql"
	"beedance-mcp/pkg/httputils"
	"beedance-mcp/pkg/loggers"
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

func MetricsServiceRelationToolSchema() mcp.Tool {
	return mcp.NewTool(
		metricsServiceRelationsToolName,
		mcp.WithDescription(metricsServiceRelationsToolDesc),
		mcp.WithString(tools.StartParamName, mcp.Description(tools.StartParamDesc)),
		mcp.WithArray(tools.ServiceNamesParamName, mcp.Required(), mcp.Description(tools.ServiceNamesParamDesc)),
	)
}

func InvokeMetricsServiceRelationTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 1. 转换请求参数
	clientVariables, err := convert2ClientVariables(request)
	if err != nil {
		loggers.Error("get client variables failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultError("获取graphql请求参数失败：" + err.Error()), nil
	}
	serverVariables, err := convert2ServerVariables(request)
	if err != nil {
		loggers.Error("get server variables failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultError("获取graphql请求参数失败：" + err.Error()), nil
	}

	// 2. 发送graphql请求
	headers, err := httputils.BuildHeaders(request)
	if err != nil {
		loggers.Error("build graphql request header failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultError("构建graphql请求头失败：" + err.Error()), nil
	}
	clientGraphqlResp, err := graphql.DoGraphqlRequest[ServiceRelationMetricsVariables, ServiceRelationClientMetricsResponse](serviceRelationClientGraphqlQuery, headers, clientVariables)
	if err != nil {
		loggers.Error("call GraphqlRequest failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultError("调用graphql接口失败：" + err.Error()), nil
	}
	serverGraphqlResp, err := graphql.DoGraphqlRequest[ServiceRelationMetricsVariables, ServiceRelationServerMetricsResponse](serviceRelationServerGraphqlQuery, headers, serverVariables)
	if err != nil {
		loggers.Error("call GraphqlRequest failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultErrorf("调用graphql接口失败：%s", err.Error()), nil
	}

	// 3. 将工具调用的结果转换成白话文
	serviceRelationClientMetrics := clientGraphqlResp.Data
	serviceRelationServerMetrics := serverGraphqlResp.Data
	loggers.Info("call graphql request success", zap.Any("request", request),
		zap.Any("service relation client metrics", serviceRelationClientMetrics),
		zap.Any("service relation server metrics", serviceRelationServerMetrics))
	message := convert2Message(request, serviceRelationClientMetrics, serviceRelationServerMetrics)
	loggers.Info("invoke service relation metrics tool success", zap.String("message", message))
	return mcp.NewToolResultText(message), nil
}
