package list_services

import (
	"beedance-mcp/api/tools/apm"
	"beedance-mcp/internal/pkg/graphql"
	"beedance-mcp/pkg/httputils"
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func ListServicesToolSchema() mcp.Tool {
	return mcp.NewTool(
		listServicesToolName,
		mcp.WithDescription(listServicesToolDesc),
		mcp.WithString(apm.WorkspaceIdParamName, mcp.Required(), mcp.Description(apm.WorkspaceIdParamDesc)),
		mcp.WithString(apm.TokenParamName, mcp.Required(), mcp.Description(apm.TokenParamDesc)),
	)
}

func InvokeListServicesTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 1. 转换请求参数
	InitServicesRegister(request)
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
	listServicesResp := graphqlResp.Data

	serviceRegister.refresh(variables.WorkspaceID, listServicesResp)

	// 3. 将工具调用结果转换成白话文
	message := convert2Message(listServicesResp)
	return mcp.NewToolResultText(message), nil
}
