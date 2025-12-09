package services_topology

import (
	"beedance-mcp/api/tools/apm"
	"beedance-mcp/internal/pkg/graphql"
	"beedance-mcp/pkg/httputils"
	"beedance-mcp/pkg/loggers"
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

func MetricsServiceToolSchema() mcp.Tool {
	return mcp.NewTool(
		serviceTopologyToolName,
		mcp.WithDescription(serviceTopologyToolDesc),
		mcp.WithString(apm.StartParamName, mcp.Description(apm.StartParamDesc)),
		mcp.WithArray(apm.ServiceNamesParamName, mcp.Required(), mcp.Description(apm.ServiceNamesParamDesc)),
	)
}

func InvokeServicesTopologyTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 1. 转换请求参数
	variables, err := convert2TopoVariables(request)
	if err != nil {
		loggers.Error("convert to topo graph request variables failed", zap.Error(err), zap.Any("request", request))
		return mcp.NewToolResultError("获取graphql请求参数失败：" + err.Error()), nil
	}

	// 2. 发送graphql请求
	headers, err := httputils.BuildHeaders(request)
	if err != nil {
		loggers.Error("build headers failed", zap.Error(err), zap.Any("request", request))
		return mcp.NewToolResultError("构建graphql请求头失败：" + err.Error()), nil
	}
	graphqlResp, err := graphql.DoGraphqlRequest[ServiceTopologyVariables, ServiceTopologyResponse](graphqlQuery, headers, variables)
	if err != nil {
		loggers.Error("send graphql request failed", zap.Error(err), zap.Any("variables", variables), zap.Any("headers", headers))
		return mcp.NewToolResultError("调用graphql接口失败：" + err.Error()), nil
	}

	// 3. 将工具调用的结果转换成白话文
	loggers.Info("call graphql request success", zap.Any("graphqlResp", graphqlResp))
	message := convert2Message(variables.WorkspaceID, graphqlResp.Data)
	loggers.Info("tool invoke success", zap.String("message", message))
	return mcp.NewToolResultText(message), nil
}
