package endpoints_traces

import (
	"beedance-mcp/api/tools"
	"beedance-mcp/api/tools/trace/list_traces"
	"beedance-mcp/internal/pkg/convertor"
	"beedance-mcp/pkg/httputils"
	"beedance-mcp/pkg/loggers"
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

func EndpointsTracesToolSchema() mcp.Tool {
	return mcp.NewTool(
		endpointsTracesToolName,
		mcp.WithDescription(endpointsTracesToolDesc),
		mcp.WithString(tools.StartParamName, mcp.Description(tools.StartParamDesc)),
		mcp.WithString(traceStateParamName, mcp.Description(traceStateParamDesc)),
		mcp.WithString(tools.ServiceNameParamName, mcp.Required(), mcp.Description(tools.ServiceNameParamDesc)),
		mcp.WithArray(endpointNamesParamName, mcp.Required(), mcp.Description(endpointNamesParamDesc)),
		mcp.WithOutputSchema[[]list_traces.ListTracesResponse](),
	)
}

func InvokeEndpointsTracesTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 1. 获取查询参数
	variables, err := convert2Variables(request)
	if err != nil {
		loggers.Error("convert request to variables failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultErrorf("获取graphql请求变量失败：%v", err), nil
	}
	headers, err := httputils.BuildHeaders(request)
	if err != nil {
		loggers.Error("build headers from request failed", zap.Any("request", request), zap.Error(err))
		mcp.NewToolResultErrorf("获取graphql请求头失败：%v", err)
	}

	// 2. 获取链路信息
	respCh := make(chan list_traces.ListTracesResponse)
	var wg sync.WaitGroup
	wg.Add(len(variables))
	for _, variable := range variables {
		go func(tracesVariable list_traces.ListTracesVariable) {
			defer wg.Done()
			endpointName := convertor.ConvertEndpointID2Name(variable.Condition.EndpointId)
			if endpointName == "" {
				loggers.Error("convert endpointId to endpointName failed", zap.Any("endpointId", variable.Condition.EndpointId), zap.Error(err))
				return
			}
			listTracesResp, err := list_traces.ListTraces(variable, headers)
			if err != nil {
				loggers.Error("list traces failed", zap.Error(err))
				return
			}
			if len(listTracesResp.Data.Traces) > 0 {
				respCh <- listTracesResp
			}
		}(variable)
	}
	go func() {
		wg.Wait()
		close(respCh)
	}()

	// 3. 采集接口的链路信息
	var toolInvokeMessageBuffer bytes.Buffer
	resps := make([]list_traces.ListTracesResponse, 0)
	for resp := range respCh {
		toolInvokeMessageBuffer.WriteString(fmt.Sprintf("接口%s的链路详情如下：\n", resp.Data.Traces[0].EndpointNames[0]))
		toolInvokeMessageBuffer.WriteString(list_traces.Convert2Message(resp.Data))
		resps = append(resps, resp)
	}
	if len(resps) == 0 {
		toolInvokeMessageBuffer.WriteString("未查询到任何链路信息")
	}
	loggers.Info("endpoints traces resp", zap.Any("response", resps))
	message := toolInvokeMessageBuffer.String()
	loggers.Info("endpoint traces tool invoke success", zap.String("message", message))
	return mcp.NewToolResultStructured(resps, message), nil
}
