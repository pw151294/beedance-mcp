package detail_trace

import (
	"beedance-mcp/internal/pkg/convertor"
	"beedance-mcp/pkg/loggers"
	"bytes"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

func convert2Variable(request mcp.CallToolRequest) (DetailTraceVariable, error) {
	traceId, err := request.RequireString(traceIDParamName)
	if err != nil {
		loggers.Error("parse traceId from request failed", zap.Any("request", request), zap.Error(err))
		return DetailTraceVariable{}, fmt.Errorf("获取链路ID参数失败：%w", err)
	}

	return DetailTraceVariable{TraceId: traceId}, nil
}

func convertTags2Message(tags []Tag) string {
	if len(tags) == 0 {
		return ""
	}

	var tagMessageBuffer bytes.Buffer
	tagMessageBuffer.WriteString("[")
	for _, tag := range tags {
		pairMessage := fmt.Sprintf("%s: %s；", tag.Key, tag.Value)
		tagMessageBuffer.WriteString(pairMessage)
	}

	tagMessage := tagMessageBuffer.String()
	return strings.TrimSuffix(tagMessage, "；") + "]"
}

func convertSpan2Message(span Span) string {
	serviceName := convertor.ConvertServiceCode2Name(span.ServiceCode)
	duration := span.EndTime - span.StartTime
	spanState := convertor.ConvertBool2Desc(span.IsError)
	tagMessage := convertTags2Message(span.Tags)
	return fmt.Sprintf(spanInfoPattern, serviceName, span.EndpointName, duration, span.Component, spanState, span.Layer, tagMessage)
}

func convert2Message(traceDetail TraceDetail) string {
	var toolInvokeMessageBuffer bytes.Buffer
	spans := traceDetail.Spans
	if len(spans) == 0 {
		toolInvokeMessageBuffer.WriteString("该链路信息为空")
	} else {
		toolInvokeMessageBuffer.WriteString("该链路信息如下：\n")
		for _, span := range spans {
			toolInvokeMessageBuffer.WriteString(convertSpan2Message(span))
		}
	}
	return toolInvokeMessageBuffer.String()
}
