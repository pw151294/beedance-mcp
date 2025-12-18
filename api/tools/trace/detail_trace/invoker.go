package detail_trace

import (
	"beedance-mcp/internal/pkg/graphql"
	"beedance-mcp/pkg/httputils"
	"beedance-mcp/pkg/loggers"
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

func DetailTraceToolSchema() mcp.Tool {
	return mcp.NewTool(
		detailTraceToolName,
		mcp.WithDescription(detailTraceToolDesc),
		mcp.WithString(traceIDParamName, mcp.Required(), mcp.Description(traceIDParamDesc)),
		mcp.WithOutputSchema[DetailTraceResponse](),
	)
}

func DetailTracesToolSchema() mcp.Tool {
	return mcp.NewTool(
		detailTracesToolName,
		mcp.WithDescription(detailTracesToolDesc),
		mcp.WithArray(traceIDsParamName, mcp.Required(), mcp.Description(traceIDsParamDesc)),
		mcp.WithOutputSchema[[]DetailTraceResponse](),
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
	return mcp.NewToolResultStructured(traceDetail, message), nil
}

func InvokeDetailTracesTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 1. 转换请求参数
	variables, err := convert2Variables(request)
	if err != nil {
		loggers.Error("convert request to graphql variables failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultErrorf("获取graphql请求变量失败：%v", err), nil
	}

	// 2. 构建graphql请求头
	headers, err := httputils.BuildHeaders(request)
	if err != nil {
		loggers.Error("build headers from request failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultError("获取graphql请求头失败：" + err.Error()), nil
	}

	// 3. 并发发送graphql请求
	respCh := make(chan DetailTraceResponse)
	var wg sync.WaitGroup
	wg.Add(len(variables))
	for _, variable := range variables {
		go func(DetailTraceVariable) {
			defer wg.Done()
			graphqlResp, err := graphql.DoGraphqlRequest[DetailTraceVariable, DetailTraceResponse](detailTraceGraphqlQuery, headers, variable)
			if err != nil {
				loggers.Error("send graphql request failed", zap.Error(err))
				return
			}
			detailTraceResp := graphqlResp.Data
			if len(detailTraceResp.TraceDetail.Spans) > 0 {
				respCh <- detailTraceResp
			}
		}(variable)
	}
	go func() {
		wg.Wait()
		close(respCh)
	}()

	// 4. 收集所有的链路信息
	var toolInvokeMessageBuffer bytes.Buffer
	var resps []DetailTraceResponse
	for resp := range respCh {
		toolInvokeMessageBuffer.WriteString(fmt.Sprintf("链路%s的详情如下：\n", resp.TraceDetail.Spans[0].TraceId))
		toolInvokeMessageBuffer.WriteString(convert2Message(resp.TraceDetail))
		resps = append(resps, resp)
	}
	if len(resps) == 0 {
		toolInvokeMessageBuffer.WriteString("未查询到任何链路信息")
	}
	loggers.Info("detail traces", zap.Any("detail traces", resps))
	message := toolInvokeMessageBuffer.String()
	loggers.Info("tool invoke success", zap.String("message", message))
	return mcp.NewToolResultStructured(resps, message), nil
}
