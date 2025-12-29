package extractor

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

// ExtractEndpointTraces 从消息体"接口GET:/permission/dataGroupTree的链路详情如下：\n链路ID：[be82591448aa43de9951c24e4ddb9fbf.71.17665653091080053]；接口：[GET:/permission/dataGroupTree]；总持续时长：40毫秒；链路状态：成功\n链路ID：[b8f689f8a98049dfaff9e569a93bd0b9.77.17667339108241687]；接口：[GET:/permission/dataGroupTree]；总持续时长：6毫秒；链路状态：成功\n"
// 提取出接口名称还有对应的链路ID数组
func ExtractEndpointTraces(text string) []EndpointTrace {
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

// ExtractSlowEndpointTraces 从消息体中提取出每个接口的慢链路ID（响应时间 > 500ms）
func ExtractSlowEndpointTraces(text string) []EndpointTrace {
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
				// 收集该接口的所有慢链路ID
				traceIds := extractSlowTraceIdsFromLines(lines[i+1:])
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

// 从"接口%s的链路详情如下：\n"提取出接口名称
func extractEndpointName(line string) string {
	suffix := strings.TrimSpace(endpointTraceDetailSuffix)
	if strings.HasSuffix(line, suffix) {
		name := strings.TrimSuffix(line, suffix)
		name = strings.TrimSpace(name)
		return strings.TrimPrefix(name, endpointNamePrefix)
	}
	return ""
}

// 从消息体"链路ID：[be82591448aa43de9951c24e4ddb9fbf.71.17665653091080053]；接口：[GET:/permission/dataGroupTree]；总持续时长：40毫秒；链路状态：成功\n链路ID：[b8f689f8a98049dfaff9e569a93bd0b9.77.17667339108241687]；接口：[GET:/permission/dataGroupTree]；总持续时长：6毫秒；链路状态：成功\n"
// 内提取出链路ID数组
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

// 从消息体行中提取出慢链路ID（响应时间 > 500ms）
func extractSlowTraceIdsFromLines(lines []string) []string {
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

		// 提取链路ID和响应时间
		if strings.HasPrefix(line, traceIdPrefix) {
			traceId, duration := ExtractTraceIdAndRt(line)
			// 只有响应时间大于阈值且traceId有效时才添加
			if traceId != "" && duration > slowTraceThreshold {
				traceIds = append(traceIds, traceId)
			}
		}
	}

	return traceIds
}

// 从"链路ID：[%s]；接口：%s；总持续时长：%d毫秒；链路状态：%s\n"提取出链路ID
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

// ExtractEndpointIdAndSla 从"服务：%s；接口：%s；接口ID：%s；成功率：%.2f\n"里提取出接口ID还有和成功率
func ExtractEndpointIdAndSla(message string) (string, int64) {
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

// ExtractEndpointIDAndRt 从"服务：%s；接口：%s；接口ID：%s；响应时间：%d毫秒\n"提取出接口ID和响应时间
func ExtractEndpointIDAndRt(message string) (string, int64) {
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

// ExtractTraceIdAndRt 从消息体"链路ID：[%s]；接口：%s；总持续时长：%d毫秒；链路状态：%s\n"提取出链路ID还有总持续时长
func ExtractTraceIdAndRt(message string) (string, int64) {
	traceId := extractTraceId(message)
	var duration int64
	var err error

	parts := strings.Split(message, semicolonSplitter)
	for _, part := range parts {
		part = strings.TrimSpace(part)

		// 提取总持续时长
		if strings.HasPrefix(part, durationLabel) {
			timePart := strings.TrimPrefix(part, durationLabel)
			timePart = strings.TrimSuffix(timePart, durationUnit)
			duration, err = strconv.ParseInt(strings.TrimSpace(timePart), 10, 64)
			if err != nil {
				loggers.Error("持续时间格式错误", zap.Error(err), zap.String("part", part))
				return "", 0
			}
		}
	}
	if traceId == "" || duration == 0 {
		loggers.Error("未找到链路ID或持续时间信息", zap.String("message", message))
		return "", 0
	}

	return traceId, duration
}
