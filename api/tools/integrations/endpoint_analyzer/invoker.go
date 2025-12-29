package endpoint_analyzer

import (
	"beedance-mcp/api/tools"
	"beedance-mcp/api/tools/trace/detail_trace"
	"beedance-mcp/api/tools/trace/list_traces"
	"beedance-mcp/internal/pkg/convertor"
	"beedance-mcp/internal/pkg/extractor"
	"beedance-mcp/pkg/loggers"
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

func EndpointErrorAnalyzerToolSchema() mcp.Tool {
	tool := mcp.NewTool(
		endpointErrorAnalyzerToolName,
		mcp.WithDescription(endpointErrorAnalyzerToolDesc),
		mcp.WithString(tools.StartParamName, mcp.Description(tools.StartParamDesc)),
		mcp.WithArray(endpointIdsParamName, mcp.Required(), mcp.Description(endpointIdsParamDesc)),
	)
	return tool
}

func EndpointSlowAnalyzerToolSchema() mcp.Tool {
	tool := mcp.NewTool(
		endpointSlowAnalyzerToolName,
		mcp.WithDescription(endpointSlowAnalyzerToolDesc),
		mcp.WithString(tools.StartParamName, mcp.Description(tools.StartParamDesc)),
		mcp.WithArray(endpointIdsParamName, mcp.Required(), mcp.Description(endpointIdsParamDesc)),
	)
	return tool
}

func InvokeEndpointErrorAnalyzerTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	service2EndpointTracesMcpRequest, err := convert2EndpointTracesMcpRequest(request, "ERROR")
	if err != nil {
		loggers.Error("convert request to endpoint traces mcp request failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultError("转换endpoint_traces工具调用参数失败：" + err.Error()), nil
	}

	// 分别调用endpoint_traces工具
	messageCh := make(chan string)
	var wg, innerWg sync.WaitGroup
	wg.Add(len(service2EndpointTracesMcpRequest))
	for serviceName, endpointsTracesMcpRequest := range service2EndpointTracesMcpRequest {
		go func(string, mcp.CallToolRequest) {
			defer wg.Done()

			endpointsTracesMcpResult, _ := list_traces.InvokeEndpointsTracesTool(ctx, endpointsTracesMcpRequest)
			endpointsTracesMcpResultText := convertor.ConvertToolCallResult2Text(endpointsTracesMcpResult)
			endpointTraces := extractor.ExtractEndpointTraces(endpointsTracesMcpResultText)
			if len(endpointTraces) == 0 {
				loggers.Info("未查询到任何错误链路", zap.String("serviceName", serviceName), zap.String("endpointsTracesMcpResultText", endpointsTracesMcpResultText))
				return
			}
			endpoint2DetailTraceMcpRequest := convert2DetailTracesMcpRequest(request, endpointTraces)

			innerWg.Add(len(endpoint2DetailTraceMcpRequest))
			for endpointName, detailTraceMcpRequest := range endpoint2DetailTraceMcpRequest {
				go func(string, mcp.CallToolRequest) {
					defer innerWg.Done()

					detailTracesMcpResult, _ := detail_trace.InvokeDetailTracesTool(ctx, detailTraceMcpRequest)
					detailTraceMcpResultText := convertor.ConvertToolCallResult2Text(detailTracesMcpResult)
					var toolInvokeMessageBuffer bytes.Buffer
					toolInvokeMessageBuffer.WriteString(fmt.Sprintf("接口%s关联服务%s，错误链路详情：\n", endpointName, serviceName))
					toolInvokeMessageBuffer.WriteString(detailTraceMcpResultText)
					messageCh <- toolInvokeMessageBuffer.String()
				}(endpointName, detailTraceMcpRequest)
			}
			innerWg.Wait()
		}(serviceName, endpointsTracesMcpRequest)
	}
	go func() {
		wg.Wait()
		close(messageCh)
	}()

	messages := make([]string, 0, 0)
	for msg := range messageCh {
		messages = append(messages, msg)
	}
	if len(messages) == 0 {
		return mcp.NewToolResultText("未查询到任何错误链路"), nil
	}
	return mcp.NewToolResultText(strings.Join(messages, "")), nil
}

func InvokeEndpointSlowAnalyzerTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	service2EndpointTracesMcpRequest, err := convert2EndpointTracesMcpRequest(request, "SUCCESS")
	if err != nil {
		loggers.Error("convert request to endpoint traces mcp request failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultError("转换endpoint_traces工具调用参数失败：" + err.Error()), nil
	}

	// 分别调用endpoint_traces工具
	messageCh := make(chan string)
	var wg, innerWg sync.WaitGroup
	wg.Add(len(service2EndpointTracesMcpRequest))
	for serviceName, endpointsTracesMcpRequest := range service2EndpointTracesMcpRequest {
		go func(string, mcp.CallToolRequest) {
			defer wg.Done()

			endpointsTracesMcpResult, _ := list_traces.InvokeEndpointsTracesTool(ctx, endpointsTracesMcpRequest)
			endpointsTracesMcpResultText := convertor.ConvertToolCallResult2Text(endpointsTracesMcpResult)
			endpointTraces := extractor.ExtractSlowEndpointTraces(endpointsTracesMcpResultText)
			if len(endpointTraces) == 0 {
				loggers.Info("未查询到任何慢链路", zap.String("serviceName", serviceName), zap.String("endpointsTracesMcpResultText", endpointsTracesMcpResultText))
				return
			}
			endpoint2DetailTraceMcpRequest := convert2DetailTracesMcpRequest(request, endpointTraces)

			innerWg.Add(len(endpoint2DetailTraceMcpRequest))
			for endpointName, detailTraceMcpRequest := range endpoint2DetailTraceMcpRequest {
				go func(string, mcp.CallToolRequest) {
					defer innerWg.Done()

					detailTracesMcpResult, _ := detail_trace.InvokeDetailTracesTool(ctx, detailTraceMcpRequest)
					detailTraceMcpResultText := convertor.ConvertToolCallResult2Text(detailTracesMcpResult)
					var toolInvokeMessageBuffer bytes.Buffer
					toolInvokeMessageBuffer.WriteString(fmt.Sprintf("接口%s关联服务%s，慢链路详情：\n", endpointName, serviceName))
					toolInvokeMessageBuffer.WriteString(detailTraceMcpResultText)
					messageCh <- toolInvokeMessageBuffer.String()
				}(endpointName, detailTraceMcpRequest)
			}
			innerWg.Wait()
		}(serviceName, endpointsTracesMcpRequest)
	}
	go func() {
		wg.Wait()
		close(messageCh)
	}()

	messages := make([]string, 0, 0)
	for msg := range messageCh {
		messages = append(messages, msg)
	}
	if len(messages) == 0 {
		return mcp.NewToolResultText("未查询到任何慢链路"), nil
	}
	return mcp.NewToolResultText(strings.Join(messages, "")), nil
}
