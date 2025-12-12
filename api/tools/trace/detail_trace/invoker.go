package detail_trace

import (
	"beedance-mcp/internal/pkg/graphql"
	"beedance-mcp/pkg/httputils"
	"beedance-mcp/pkg/loggers"
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

func DetailTraceToolSchema() mcp.Tool {
	return mcp.NewTool(
		detailTraceToolName,
		mcp.WithDescription(detailTraceToolDesc),
		mcp.WithString(traceIDParamName, mcp.Required(), mcp.Description(traceIDParamDesc)),
	)
}

func InvokeDetailTraceTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 1. 转换请求参数
	variable, err := convert2Variable(request)
	if err != nil {
		loggers.Error("convert tool call request to graphql variable failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultError("获取graphql请求变量失败：" + err.Error()), nil
	}

	// 2. 发送graphql请求
	headers, err := httputils.BuildHeaders(request)
	if err != nil {
		loggers.Error("build headers from request failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultError("获取graphql请求头失败：" + err.Error()), nil
	}
	graphqlResp, err := graphql.DoGraphqlRequest[DetailTraceVariable, DetailTraceResponse](detailTraceGraphqlQuery, headers, variable)
	if err != nil {
		loggers.Error("send graphql request failed", zap.Error(err))
		return mcp.NewToolResultError("graphql 请求失败：" + err.Error()), nil
	}

	// 3. 将工具调用的结果转换成白话文
	traceDetail := graphqlResp.Data
	loggers.Info("trace detail", zap.Any("trace detail", traceDetail))
	message := convert2Message(traceDetail.TraceDetail)
	loggers.Info("tool invoke success", zap.String("tool invoke message", message))
	return mcp.NewToolResultText(message), nil
}
