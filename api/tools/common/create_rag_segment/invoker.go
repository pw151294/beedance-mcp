package create_rag_segment

import (
	"beedance-mcp/configs"
	"beedance-mcp/internal/pkg/graphql"
	"beedance-mcp/pkg/httputils"
	"beedance-mcp/pkg/loggers"
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

func InvokeCreateSegmentToolSchema() mcp.Tool {
	return mcp.NewTool(
		ragSegmentCreateToolName,
		mcp.WithDescription(ragSegmentCreateToolDesc),
		mcp.WithString(ragDatasetIdParamName, mcp.Required(), mcp.Description(ragDatasetIdParamDesc)),
		mcp.WithString(ragDocumentIdParamName, mcp.Required(), mcp.Description(ragDocumentIdParamDesc)),
		mcp.WithString(contentParamName, mcp.Required(), mcp.Description(contentParamDesc)),
	)
}

func InvokeCreateRagSegmentTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 1. 获取查询参数
	variable, err := convert2Variable(request)
	if err != nil {
		loggers.Error("convert request to segment create variable failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultError("获取HTTP请求变量失败：" + err.Error()), nil
	}
	headers, err := httputils.BuildHeaders(request)
	if err != nil {
		loggers.Error("build headers from request failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultErrorf("获取HTTP请求头失败：%v", err), nil
	}

	// 2. 执行工具调用
	url := configs.GlobalConfig.Gateway.BeedanceAddress + ragCreateSegmentUrl
	httpResp, err := graphql.DoHttpRequest[RagSegmentCreateVariable, RagSegmentCreateResponse](url, headers, variable)
	if err != nil {
		loggers.Error("send segment create http request failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultError("发送HTTP请求失败：" + err.Error()), nil
	}
	if httpResp.Success {
		return mcp.NewToolResultText("分析报告已入库"), nil
	} else {
		return mcp.NewToolResultError("分析报告入库失败：" + httpResp.Message), nil
	}
}
