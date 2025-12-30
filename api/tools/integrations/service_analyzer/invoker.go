package service_analyzer

import (
	"beedance-mcp/api/tools"
	"beedance-mcp/api/tools/apm/metrics_endpoints"
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

func ServiceErrorAnalyzerToolSchema() mcp.Tool {
	return mcp.NewTool(
		serviceErrorAnalyzerToolName,
		mcp.WithDescription(serviceErrorAnalyzerToolDesc),
		mcp.WithString(tools.ServiceNameParamName, mcp.Required(), mcp.Description(tools.ServiceNameParamDesc)),
		mcp.WithString(tools.StartParamName, mcp.Description(tools.StartParamDesc)))
}

func ServiceSlowAnalyzerToolSchema() mcp.Tool {
	return mcp.NewTool(
		serviceSlowAnalyzerToolName,
		mcp.WithDescription(serviceSlowAnalyzerToolDesc),
		mcp.WithString(tools.ServiceNameParamName, mcp.Required(), mcp.Description(tools.ServiceNameParamDesc)),
		mcp.WithString(tools.StartParamName, mcp.Description(tools.StartParamDesc)))
}

func InvokeServiceErrorAnalyzerTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 1. 调用metrics_endpoints工具查询该服务的所有接口和成功率
	metricsEndpointsMcpRequest, err := convert2MetricsEndpointsMcpRequest(request, endpointSlaMetricsName)
	if err != nil {
		loggers.Error("convert request to metrics endpoints mcp request failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultErrorf("构建metrics_endpoints工具调用参数失败：%v", err), nil
	}
	metricsEndpointsMcpResult, err := metrics_endpoints.InvokeMetricsEndpointsTool(ctx, metricsEndpointsMcpRequest)
	if err != nil {
		loggers.Error("call metrics endpoints mcp request failed", zap.Any("request", metricsEndpointsMcpRequest), zap.Error(err))
		return mcp.NewToolResultErrorf("调用metrics_endpoints工具失败：%v", err), nil
	}
	metricsEndpointsMcpResultText := convertor.ConvertToolCallResult2Text(metricsEndpointsMcpResult)

	// 2. 从metrics_endpoints工具调用的结果里筛选出异常的接口
	var endpointIds []string
	if strings.HasPrefix(metricsEndpointsMcpResultText, metricsEndpointMessagePrefix) {
		messages := strings.Split(metricsEndpointsMcpResultText, lineBreaker)
		messages = messages[1:]
		for _, message := range messages {
			if message == "" {
				continue
			}
			endpointId, sla := extractor.ExtractEndpointIdAndSla(message)
			if sla < 10000 {
				endpointIds = append(endpointIds, endpointId)
			}
		}
	}
	if len(endpointIds) == 0 {
		return mcp.NewToolResultText("该服务下没有错误链路"), nil
	}

	// 3. 根据endpointIds调用endpoints_traces工具
	endpointsTracesMcpRequest := convert2EndpointTracesMcpRequest(request, endpointIds, endpointSlaMetricsName)
	endpointsTracesMcpResult, _ := list_traces.InvokeEndpointsTracesTool(ctx, endpointsTracesMcpRequest)
	endpointsTracesMcpResultText := convertor.ConvertToolCallResult2Text(endpointsTracesMcpResult)
	endpointTraces := extractor.ExtractEndpointTraces(endpointsTracesMcpResultText)
	if len(endpointTraces) == 0 {
		return mcp.NewToolResultText("该服务下没有错误链路"), nil
	}

	// 4. 根据endpointName还有traceIds调用detail_traces工具
	endpoint2DetailTraceMcpRequest := convert2DetailTracesMcpRequest(request, endpointTraces)
	var wg sync.WaitGroup
	wg.Add(len(endpointTraces))
	messageCh := make(chan string)
	for endpointName, detailTraceMcpRequest := range endpoint2DetailTraceMcpRequest {
		go func(string, mcp.CallToolRequest) {
			defer wg.Done()

			detailTracesMcpResult, _ := detail_trace.InvokeDetailTracesTool(ctx, detailTraceMcpRequest)
			detailTraceMcpResultText := convertor.ConvertToolCallResult2Text(detailTracesMcpResult)
			var toolInvokeMessageBuffer bytes.Buffer
			toolInvokeMessageBuffer.WriteString(fmt.Sprintf("接口%s的错误链路详情如下：\n", endpointName))
			toolInvokeMessageBuffer.WriteString(detailTraceMcpResultText)
			messageCh <- toolInvokeMessageBuffer.String()
		}(endpointName, detailTraceMcpRequest)
	}
	go func() {
		wg.Wait()
		close(messageCh)
	}()

	// 5. 采集各接口的错误链路详情
	serviceName, _ := request.RequireString(tools.ServiceNameParamName)
	messages := []string{fmt.Sprintf("服务%s的错误链路详情如下：\n", serviceName)}
	for msg := range messageCh {
		messages = append(messages, msg)
	}
	return mcp.NewToolResultText(strings.Join(messages, lineBreaker)), nil
}

func InvokeServiceSlowAnalyzerTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 1. 调用metrics_endpoints工具查询该服务的所有接口的响应时间
	metricsEndpointsMcpRequest, err := convert2MetricsEndpointsMcpRequest(request, endpointRtMetricsName)
	if err != nil {
		loggers.Error("convert request to metrics endpoints mcp request failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultErrorf("构建metrics_endpoints工具调用参数失败：%v", err), nil
	}
	metricsEndpointsMcpResult, err := metrics_endpoints.InvokeMetricsEndpointsTool(ctx, metricsEndpointsMcpRequest)
	if err != nil {
		loggers.Error("call metrics endpoints mcp request failed", zap.Any("request", metricsEndpointsMcpRequest), zap.Error(err))
		return mcp.NewToolResultErrorf("调用metrics_endpoints工具失败：%v", err), nil
	}
	metricsEndpointsMcpResultText := convertor.ConvertToolCallResult2Text(metricsEndpointsMcpResult)

	// 2. 从metrics_endpoints工具调用的结果里筛选出响应时间长的接口
	var endpointIds []string
	if strings.HasPrefix(metricsEndpointsMcpResultText, metricsEndpointMessagePrefix) {
		messages := strings.Split(metricsEndpointsMcpResultText, lineBreaker)
		messages = messages[1:]
		for _, message := range messages {
			if message == "" {
				continue
			}
			endpointId, rt := extractor.ExtractEndpointIDAndRt(message)
			if rt > 200 {
				endpointIds = append(endpointIds, endpointId)
			}
		}
	}
	if len(endpointIds) == 0 {
		return mcp.NewToolResultText("该服务下没有慢链路"), nil
	}

	// 3. 根据endpointIds调用endpoints_traces工具
	endpointsTracesMcpRequest := convert2EndpointTracesMcpRequest(request, endpointIds, endpointRtMetricsName)
	endpointsTracesMcpResult, _ := list_traces.InvokeEndpointsTracesTool(ctx, endpointsTracesMcpRequest)
	endpointsTracesMcpResultText := convertor.ConvertToolCallResult2Text(endpointsTracesMcpResult)
	endpointTraces := extractor.ExtractEndpointTraces(endpointsTracesMcpResultText)
	if len(endpointTraces) == 0 {
		return mcp.NewToolResultText("该服务下没有慢链路"), nil
	}

	// 4. 根据endpointName还有traceIds调用detail_traces工具
	endpoint2DetailTraceMcpRequest := convert2DetailTracesMcpRequest(request, endpointTraces)
	var wg sync.WaitGroup
	wg.Add(len(endpointTraces))
	messageCh := make(chan string)
	for endpointName, detailTraceMcpRequest := range endpoint2DetailTraceMcpRequest {
		go func(string, mcp.CallToolRequest) {
			defer wg.Done()

			detailTracesMcpResult, _ := detail_trace.InvokeDetailTracesTool(ctx, detailTraceMcpRequest)
			detailTraceMcpResultText := convertor.ConvertToolCallResult2Text(detailTracesMcpResult)
			var toolInvokeMessageBuffer bytes.Buffer
			toolInvokeMessageBuffer.WriteString(fmt.Sprintf("接口%s的慢链路详情如下：\n", endpointName))
			toolInvokeMessageBuffer.WriteString(detailTraceMcpResultText)
			messageCh <- toolInvokeMessageBuffer.String()
		}(endpointName, detailTraceMcpRequest)
	}
	go func() {
		wg.Wait()
		close(messageCh)
	}()

	// 5. 采集各接口的错误链路详情
	serviceName, _ := request.RequireString(tools.ServiceNameParamName)
	messages := []string{fmt.Sprintf("服务%s的慢链路详情如下：\n", serviceName)}
	for msg := range messageCh {
		messages = append(messages, msg)
	}
	return mcp.NewToolResultText(strings.Join(messages, lineBreaker)), nil

}
