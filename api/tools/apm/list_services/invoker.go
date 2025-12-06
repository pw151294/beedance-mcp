package list_services

import (
	"beedance-mcp/internal/pkg/graphql"
	"beedance-mcp/pkg/httputils"
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func ListServicesToolSchema() mcp.Tool {
	return mcp.NewTool(
		listServicesToolName,
		mcp.WithDescription(listServicesToolDesc),
		mcp.WithString(workspaceIdParamName, mcp.Required(), mcp.Description(workspaceIdParamDesc)),
		mcp.WithString(tokenParamName, mcp.Required(), mcp.Description(tokenParamDesc)),
	)
}

func InvokeListServicesTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 1. 转换请求参数
	variables, err := convert2Variables(request)
	if err != nil {
		return mcp.NewToolResultError("获取graphql请求参数失败：" + err.Error()), nil
	}

	// 2. 发送graphql请求
	headers, err := httputils.BuildHeaders(request)
	if err != nil {
		return mcp.NewToolResultError("构建graphql请求头失败：" + err.Error()), nil
	}
	graphqlResp, err := graphql.DoGraphqlRequest[ListServicesVariables, ListServicesResponse](graphQuery, headers, variables)
	if err != nil {
		return mcp.NewToolResultError("调用graphql接口失败：" + err.Error()), nil
	}

	// 3. 将工具调用结果转换成白话文
	message := convert2Message(graphqlResp.Data)
	return mcp.NewToolResultText(message), nil
}
