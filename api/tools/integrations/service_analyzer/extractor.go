package service_analyzer

import (
	"beedance-mcp/pkg/loggers"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type EndpointTrace struct {
	EndpointName string
	TraceIds     []string
}

func extractEndpointTraces(text string) []EndpointTrace {
	// 1. 修改入口判断条件，使用常量
	if !strings.Contains(text, endpointTraceDetailSuffix) {
		return []EndpointTrace{}
	}

	result := make([]EndpointTrace, 0)
	lines := strings.Split(text, "\n")
	suffix := strings.TrimSpace(endpointTraceDetailSuffix)

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		if strings.HasSuffix(line, suffix) {
			// 提取接口名称
			endpointName := extractEndpointName(line)
			if endpointName != "" {
				// 收集该接口的所有链路ID
				traceIds := extractTraceIdsFromLines(lines[i+1:])
				if len(traceIds) > 0 {
					result = append(result, EndpointTrace{
						EndpointName: endpointName,
						TraceIds:     traceIds,
					})
				}
			}
		}
	}

	return result
}

func extractEndpointName(line string) string {
	suffix := strings.TrimSpace(endpointTraceDetailSuffix)
	if strings.HasSuffix(line, suffix) {
		name := strings.TrimSuffix(line, suffix)
		name = strings.TrimSpace(name)
		return strings.TrimPrefix(name, endpointNamePrefix)
	}
	return ""
}

func extractTraceIdsFromLines(lines []string) []string {
	var traceIds []string
	suffix := strings.TrimSpace(endpointTraceDetailSuffix)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 如果遇到新的接口描述，停止收集
		if strings.HasSuffix(line, suffix) {
			break
		}

		// 提取链路ID，使用 HasPrefix 判断
		if strings.HasPrefix(line, traceIdPrefix) {
			traceId := extractTraceId(line)
			if traceId != "" {
				traceIds = append(traceIds, traceId)
			}
		}
	}

	return traceIds
}

func extractTraceId(line string) string {
	// 去除前缀 "链路ID：["
	if strings.HasPrefix(line, traceIdPrefix) {
		val := strings.TrimPrefix(line, traceIdPrefix)
		// 查找右括号 "]"
		if idx := strings.Index(val, rightBracket); idx != -1 {
			return val[:idx]
		}
	}
	return ""
}

func extractEndpointIdAndSla(message string) (string, int64) {
	// 分割字符串获取各个字段
	parts := strings.Split(message, semicolonSplitter)

	var endpointID string
	var slaStr string

	// 遍历分割后的字段
	for _, part := range parts {
		part = strings.TrimSpace(part)

		if strings.HasPrefix(part, endpointIdPrefix) {
			// 提取接口ID
			endpointID = strings.TrimPrefix(part, endpointIdPrefix)
		} else if strings.HasPrefix(part, endpointSlaPrefix) {
			// 提取成功率字符串
			slaStr = strings.TrimPrefix(part, endpointSlaPrefix)
		}
	}

	// 将百分比成功率转换为整数（乘以100）
	var sla int64
	if slaStr != "" {
		// 去掉百分号并转换为float64
		cleanSlaStr := strings.TrimSuffix(slaStr, "%")
		slaFloat, err := strconv.ParseFloat(cleanSlaStr, 64)
		if err != nil {
			loggers.Error("convert slaStr to float64 failed", zap.String("slaStr", cleanSlaStr), zap.Error(err))
			return endpointID, 0
		}
		// 乘以100转换为整数（99.85 -> 9985）
		sla = int64(slaFloat * 100)
	}

	return endpointID, sla
}

func extractEndpointIDAndRt(message string) (string, int64) {
	// 异常处理：空字符串
	if message == "" {
		return "", 0
	}

	var endpointID string
	var responseTime int64

	// 按分号分割字符串
	parts := strings.Split(message, semicolonSplitter)

	for _, part := range parts {
		part = strings.TrimSpace(part)

		if strings.HasPrefix(part, endpointIdPrefix) {
			endpointID = strings.TrimPrefix(part, endpointIdPrefix)
		} else if strings.HasPrefix(part, endpointRtPrefix) {
			rtStr := strings.TrimPrefix(part, endpointRtPrefix)
			rtStr = strings.TrimSuffix(rtStr, "毫秒")
			rtStr = strings.TrimSuffix(rtStr, "\n")
			rtStr = strings.TrimSpace(rtStr)

			// 转换为int64
			if rt, err := strconv.ParseInt(rtStr, 10, 64); err == nil {
				responseTime = rt
			}
		}
	}

	return endpointID, responseTime
}
