package list_services

import (
	"beedance-mcp/internal/pkg/graphql"
	"beedance-mcp/pkg/httputils"
	"beedance-mcp/pkg/loggers"
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

func ListServicesToolSchema() mcp.Tool {
	return mcp.NewTool(
		listServicesToolName,
		mcp.WithDescription(listServicesToolDesc),
	)
}

func listServices(request mcp.CallToolRequest) (ListServicesResponse, error) {
	// 1. 转换请求参数
	variables, err := convert2Variables(request)
	if err != nil {
		loggers.Error("convert request to graphql variables failed", zap.Any("request", request), zap.Error(err))
		return ListServicesResponse{}, err
	}

	// 2. 发送graphql请求
	headers, err := httputils.BuildHeaders(request)
	if err != nil {
		loggers.Error("build headers failed", zap.Any("request", request), zap.Error(err))
		return ListServicesResponse{}, err
	}
	graphqlResp, err := graphql.DoGraphqlRequest[ListServicesVariables, ListServicesResponse](graphqlQuery, headers, variables)
	if err != nil {
		loggers.Error("call graphql request failed", zap.Any("variables", variables), zap.Any("headers", headers), zap.Error(err))
		return ListServicesResponse{}, err
	}

	loggers.Info("call graphql request success", zap.Any("list services", graphqlResp))
	return graphqlResp.Data, nil
}

func InvokeListServicesTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 获取服务列表数据
	listServicesResp, err := listServices(request)
	if err != nil {
		return mcp.NewToolResultError("获取服务列表失败：" + err.Error()), nil
	}

	// 将工具调用结果转换成白话文
	message := convert2Message(listServicesResp)
	loggers.Info("tool invoke success", zap.String("message", message))
	return mcp.NewToolResultText(message), nil
}
