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

func convert2Variables(request mcp.CallToolRequest) ([]DetailTraceVariable, error) {
	traceIds, err := request.RequireStringSlice(traceIDsParamName)
	if err != nil || len(traceIds) == 0 {
		loggers.Error("parse traceIds from request failed", zap.Any("request", request), zap.Error(err))
		return nil, fmt.Errorf("获取链路ID数组失败：%w", err)
	}

	variables := make([]DetailTraceVariable, 0, len(traceIds))
	for _, traceId := range traceIds {
		variables = append(variables, DetailTraceVariable{TraceId: traceId})
	}
	return variables, nil
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

func convertLogs2Message(logs []Log) string {
	if len(logs) == 0 {
		return ""
	}

	errorLogs := make([]string, 0, 0)
	for _, log := range logs {
		data := log.Data
		if len(data) == 0 {
			continue
		}
		logTags := make(map[string]string)
		for _, tag := range data {
			logTags[tag.Key] = tag.Value
		}
		if logTags[eventPropertyName] != "error" {
			continue
		}
		message := logTags[messagePropertyName]
		stack := logTags[stackPropertyName]
		errorKind := logTags[errorKindPropertyName]
		errorLogs = append(errorLogs, fmt.Sprintf(errorLogInfoPattern, len(errorLogs)+1, errorKind, message, stack[:stackLengthThreshold]))
	}

	return strings.Join(errorLogs, "；")
}

func convertSpan2Message(span Span) string {
	serviceName := convertor.ConvertServiceCode2Name(span.ServiceCode)
	duration := span.EndTime - span.StartTime
	spanState := convertor.ConvertBool2Desc(span.IsError)
	tagMessage := convertTags2Message(span.Tags)
	logMessage := convertLogs2Message(span.Logs)
	return fmt.Sprintf(spanInfoPattern, serviceName, span.EndpointName, duration, span.Component, spanState, span.Layer, tagMessage, logMessage)
}

func convert2Message(traceDetail TraceDetail) string {
	var toolInvokeMessageBuffer bytes.Buffer
	spans := traceDetail.Spans
	if len(spans) == 0 {
		toolInvokeMessageBuffer.WriteString("该链路信息为空")
	} else {
		for _, span := range spans {
			toolInvokeMessageBuffer.WriteString(convertSpan2Message(span))
		}
	}
	return toolInvokeMessageBuffer.String()
}
