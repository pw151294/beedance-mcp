package list_traces

import (
	"beedance-mcp/api/tools"
	"beedance-mcp/internal/pkg/convertor"
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

// ListTracesToolSchema --- list_traces ---
func ListTracesToolSchema() mcp.Tool {
	return mcp.NewTool(
		listTracesToolName,
		mcp.WithDescription(listTracesToolDesc),
		mcp.WithString(tools.StartParamName, mcp.Description(tools.StartParamDesc)),
		mcp.WithString(traceStateParamName, mcp.Description(traceStateParamDesc)),
		mcp.WithString(tools.ServiceNameParamName, mcp.Required(), mcp.Description(tools.ServiceNameParamDesc)),
		mcp.WithString(endpointNameParamName, mcp.Description(endpointNameParamDesc)),
	)
}

func ListTraces(variable ListTracesVariable, headers map[string]string) (ListTracesResponse, error) {
	// 1. 发送graphql请求
	graphqlResp, err := graphql.DoGraphqlRequest[ListTracesVariable, ListTracesResponse](listTracesGraphqlQuery, headers, variable)
	if err != nil {
		loggers.Error("send graphql request failed", zap.Any("variable", variable), zap.Error(err))
		return ListTracesResponse{}, err
	}
	return graphqlResp.Data, nil
}

func InvokeListTracesTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 1. 转换请求参数
	variable, err := convert2Variable(request)
	if err != nil {
		loggers.Error("convert request to variable failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultError("获取graphql请求参数失败：" + err.Error()), nil
	}

	// 2. 发送graphql请求
	headers, err := httputils.BuildHeaders(request)
	if err != nil {
		loggers.Error("build headers from request failed", zap.Any("request", request), zap.Error(err))
		return mcp.NewToolResultError("获取graphql请求头失败：" + err.Error()), nil
	}
	listTracesResp, err := ListTraces(variable, headers)
	if err != nil {
		loggers.Error("invoke list traces failed", zap.Error(err))
		return mcp.NewToolResultError("调用list_traces工具失败：" + err.Error()), nil
	}

	// 3. 使用结构化输出
	loggers.Info("list traces response", zap.Any("list traces response", listTracesResp))
	message := convert2Message(listTracesResp.Data)
	loggers.Info("tool invoke success", zap.Any("message", message))
	return mcp.NewToolResultText(message), nil
}

// EndpointsTracesToolSchema --- endpoints_traces
func EndpointsTracesToolSchema() mcp.Tool {
	return mcp.NewTool(
		endpointsTracesToolName,
		mcp.WithDescription(endpointsTracesToolDesc),
		mcp.WithString(tools.StartParamName, mcp.Description(tools.StartParamDesc)),
		mcp.WithString(traceStateParamName, mcp.Description(traceStateParamDesc)),
		mcp.WithString(tools.ServiceNameParamName, mcp.Required(), mcp.Description(tools.ServiceNameParamDesc)),
		mcp.WithArray(endpointIdsParamName, mcp.Required(), mcp.Description(endpointIdsParamDesc)),
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
		return mcp.NewToolResultErrorf("获取graphql请求头失败：%v", err), nil
	}

	// 2. 获取链路信息
	respCh := make(chan ListTracesResponse)
	var wg sync.WaitGroup
	wg.Add(len(variables))
	for _, variable := range variables {
		go func(tracesVariable ListTracesVariable) {
			defer wg.Done()

			endpointName := convertor.ConvertEndpointID2Name(tracesVariable.Condition.EndpointId)
			if endpointName == "" {
				loggers.Error("convert endpointId to endpointName failed", zap.Any("endpointId", tracesVariable.Condition.EndpointId))
				return
			}

			listTracesResp, err := ListTraces(tracesVariable, headers)
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
	resps := make([]ListTracesResponse, 0)
	for resp := range respCh {
		toolInvokeMessageBuffer.WriteString(fmt.Sprintf("接口%s的链路详情如下：\n", resp.Data.Traces[0].EndpointNames[0]))
		toolInvokeMessageBuffer.WriteString(convert2Message(resp.Data))
		resps = append(resps, resp)
	}
	if len(resps) == 0 {
		toolInvokeMessageBuffer.WriteString("未查询到任何链路信息")
	}
	loggers.Info("endpoints traces resp", zap.Any("response", resps))
	message := toolInvokeMessageBuffer.String()
	loggers.Info("endpoint traces tool invoke success", zap.String("message", message))
	return mcp.NewToolResultText(message), nil
}
